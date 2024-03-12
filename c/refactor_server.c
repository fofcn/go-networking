#include <stdio.h>
#include <pthread.h>
#include <sys/epoll.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <signal.h>

#define PORT 12345
#define MAX_EVENTS 10
#define BUFF_SIZE 1024
#define WORKER_SIZE 4

// epoll file descriptor
int epoll_fd;

// handlers
void* handle_accept(void* arg);
void* handle_io(void* arg);
void wait_to_death();

int main() {
    int server_fd;
    struct sockaddr_in server_addr;

    // create server socket
    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    // set to non-blocking
    fcntl(server_fd, F_SETFL, fcntl(server_fd, F_GETFL, 0) | O_NONBLOCK);

    // bind
    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = htonl(INADDR_ANY);
    server_addr.sin_port = htons(PORT);

    printf("binding\n");
    bind(server_fd, (struct sockaddr *)&server_addr, sizeof(server_addr));

    printf("listen\n");
    // listen
    listen(server_fd, MAX_EVENTS);

    printf("epoll create\n");
    // create epoll and add server_fd
    epoll_fd = epoll_create1(0);

    struct epoll_event event;
    event.events = EPOLLIN;
    event.data.fd = server_fd;
    printf("epoll add\n");
    epoll_ctl(epoll_fd, EPOLL_CTL_ADD, server_fd, &event);

    // create accept handler threads
    pthread_t accept_threads[2];
    for (int i = 0; i < 2; i++) {
        printf("create acceptor thread. index: %d\n", i);
		pthread_create(&accept_threads[i], NULL, handle_accept, &server_fd);
        pthread_detach(accept_threads[i]);
	}
    
    wait_to_death();

    close(epoll_fd);
    close(server_fd);

    return 0;
}

void* handle_accept(void* arg) {
    int server_fd = *(int*)arg;

    while (1) {
		struct epoll_event events[MAX_EVENTS];
		int n = epoll_wait(epoll_fd, events, MAX_EVENTS, -1);
		for (int i = 0; i < n; i++) {
			if (events[i].data.fd == server_fd) {
				// new connection arrives
				int client_fd = accept(server_fd, NULL, NULL);

				// set to non-blocking
				fcntl(client_fd, F_SETFL, fcntl(client_fd, F_GETFL, 0) | O_NONBLOCK);
                
				// create a worker thread to handle this connection
				pthread_t worker_thread;
				pthread_create(&worker_thread, NULL, handle_io, &client_fd);
			}
		}
	}
}

void* handle_io(void* arg) {
    int client_fd = *(int*)arg;

    while (1) {
		char buff[BUFF_SIZE] = {0};
		int len = read(client_fd, buff, BUFF_SIZE);
		if (len <= 0) {
			// error occurs or the client closes connection
			close(client_fd);

			struct epoll_event event;
			event.events = EPOLLIN;
			event.data.fd = client_fd;
			epoll_ctl(epoll_fd, EPOLL_CTL_DEL, client_fd, &event);

			break;
		}
		else {
			printf("Received %s from client\n", buff);
		}
	}

	return NULL;
}

void wait_to_death() {
    sigset_t allset;
    sigemptyset(&allset);
    sigaddset(&allset, SIGINT); // Ctrl+C
    sigaddset(&allset, SIGQUIT); // Ctrl+\

    int sig;
    for (;;) {
        int err = sigwait(&allset, &sig);
        if (err == 0) {
            printf("received signal %d, prepare to exit\n", sig);
            break;
        }
    }
}
