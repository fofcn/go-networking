#include <stdio.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <string.h>
#include <unistd.h>

#define PORT 12345
#define BUF_SIZE 1024

int main() {
    int client_fd;
    struct sockaddr_in servaddr;
    char buf[BUF_SIZE] = "Hello, Server!";

    /* 创建 TCP socket */
    client_fd = socket(AF_INET, SOCK_STREAM, 0);

    /* 连接到服务器 */
    memset(&servaddr, 0, sizeof(servaddr));
    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = htonl(INADDR_ANY);
    servaddr.sin_port = htons(PORT);
    connect(client_fd, (struct sockaddr *)&servaddr, sizeof(servaddr));

    /* 发送数据 */
    send(client_fd, buf, strlen(buf)+1, 0);

    close(client_fd);
    return 0;
}