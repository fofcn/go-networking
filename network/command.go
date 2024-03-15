package network

type CommandType uint32

// CmdType 定义了连接中使用的命令类型。

const (
	CONN            CommandType = iota + 1 // 用于开始端到端的连接。客户端向服务器发送其 DH 公钥，服务器执行同样的操作。
	CONNACK                                // 对于CONN的响应。服务器向客户端发送其 DH 公钥。
	PING                                   // 用于维持会话连接。客户端向服务器发送，服务器响应PONG。
	PONG                                   // 对于PING的响应。
	CLOSE                                  // 用于关闭连接。客户端或服务器均可发起。
	CLOSEACK                               // 对于CLOSE的响应。服务器完成发送所有待发帧后发送。
	LISTDIR                                // 用于请求服务器列出目录和文件。
	LISTDIRACK                             // 对于LISTDIR的响应，包含文件列表。
	FILETRANSFER                           // 客户端向服务器发送，以协商文件传输。
	FILETRANSFERACK                        // 对于FILETRANSFER的响应，包含文件传输细节。
	TRANSFER                               // 用于实际传输文件数据。
)
