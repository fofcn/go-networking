#include <stdio.h>
#include <sys/epoll.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>

#define PORT 12345
#define BUF_SIZE 1024
#define MAX_EVENTS 10

static int setnonblocking(int sockfd);

int main() {
    int server_fd, client_fd;
    struct sockaddr_in servaddr;
    char buf[BUF_SIZE];

    /* 创建 TCP socket */
    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    
    /* 绑定到某个端口 */
    memset(&servaddr, 0, sizeof(servaddr));
    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = htonl(INADDR_ANY);
    servaddr.sin_port = htons(PORT);
    bind(server_fd, (struct sockaddr *)&servaddr, sizeof(servaddr));

    /* 开始监听 */
    listen(server_fd, MAX_EVENTS);

    /* 创建epoll并添加监听事件 */
    int epoll_fd = epoll_create1(0);
    struct epoll_event event;
    event.events = EPOLLIN;
    event.data.fd = server_fd;
    epoll_ctl(epoll_fd, EPOLL_CTL_ADD, server_fd, &event);

    /* 开始处理事件 */
    while (1) {
        struct epoll_event events[MAX_EVENTS];
        int n = epoll_wait(epoll_fd, events, 10, -1);
        for (int i = 0; i < n; i++) {
            if (events[i].data.fd == server_fd) {
                /* 新的连接到来 */
                client_fd = accept(server_fd, NULL, NULL);
                // 设置成non-blocking I/O
                setnonblocking(client_fd);
                event.data.fd = client_fd;
                // 添加这个文件描述符到epoll中
                epoll_ctl(epoll_fd, EPOLL_CTL_ADD, client_fd, &event);
            } else if (events[i].events & EPOLLIN) {
                /* 处理epoll事件 */
                for (;;) {
                    bzero(buf, sizeof(buf));
                    /* 读取数据并处理 */
                    n = read(events[i].data.fd, buf, BUF_SIZE);
                    if (n == 0) {
                        printf("client connection closed. n == 0\n");
                        close(events[i].data.fd);
                        epoll_ctl(epoll_fd, EPOLL_CTL_DEL, events[i].data.fd, NULL);
                       
                    } else if (n < 0) {
                        // 暂时读取不到数据就可以等下次epoll_wait调用
                         if(errno == EAGAIN || errno == EWOULDBLOCK) {
                            break;
                         } else { // 发生了错误，就需要关闭连接并把连接从epoll监控链中删除
                            printf("client connection closed. n < 0, actual value of n: %d\n", n);
                            close(events[i].data.fd);
                            epoll_ctl(epoll_fd, EPOLL_CTL_DEL, events[i].data.fd, NULL);       
                            break; 
                         }
                    } else {
                        printf("Received: %s\n", buf);
                    }
                }
            } else {
                printf("[+] unexpected\n");
            }

            /* 检查连接是否关闭 */
			if (events[i].events & (EPOLLRDHUP | EPOLLHUP)) {
				printf("[+] connection closed\n");
				epoll_ctl(epoll_fd, EPOLL_CTL_DEL,
					  events[i].data.fd, NULL);
				close(events[i].data.fd);
				continue;
			}
		}
    }

    close(server_fd);
    return 0;
}

static int setnonblocking(int sockfd) {
	if (fcntl(sockfd, F_SETFL, fcntl(sockfd, F_GETFL, 0) | O_NONBLOCK) ==
	    -1) {
		return -1;
	}
	return 0;
}