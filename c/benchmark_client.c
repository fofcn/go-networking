#include <stdio.h>
#include <pthread.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <unistd.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <time.h>

#define PORT 12345
#define NUM_CONNECTIONS 100 // 多少并发连接
#define REQUESTS_PER_CONN 10 // 每个连接发送多少个请求

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


void *client_thread(void *arg) {
    int sock_fd = socket(AF_INET, SOCK_STREAM, 0);
    struct sockaddr_in server_addr;

    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(PORT);
    server_addr.sin_addr.s_addr = inet_addr("127.0.0.1");

    // Connect to the server
    connect(sock_fd, (struct sockaddr *)&server_addr, sizeof(server_addr));

    char buffer[1024];

    for (int i = 0; i < REQUESTS_PER_CONN; i++) {
        // 发送数据到服务器
        sprintf(buffer, "Request #%d", i);
        send(sock_fd, buffer, strlen(buffer), 0);

        // 接收服务器的响应
        recv(sock_fd, buffer, sizeof(buffer), 0);
    }

    close(sock_fd);
    return NULL;
}

int main() {
    pthread_t threads[NUM_CONNECTIONS];

    clock_t start = clock();
    for (int i = 0; i < NUM_CONNECTIONS; i++) {
        pthread_create(&threads[i], NULL, client_thread, NULL);
    }

    for (int i = 0; i < NUM_CONNECTIONS; i++) {
        // pthread_join(threads[i], NULL);
    }
    clock_t end = clock();

    double elapsed_time = (double)(end - start) / CLOCKS_PER_SEC;

    printf("Total time taken: %f\n", elapsed_time);
    printf("Average time per request: %f\n", elapsed_time / (NUM_CONNECTIONS * REQUESTS_PER_CONN));
    printf("Throughput (requests per second): %f\n", (NUM_CONNECTIONS * REQUESTS_PER_CONN) / elapsed_time);
    wait_to_death();
    return 0;
}