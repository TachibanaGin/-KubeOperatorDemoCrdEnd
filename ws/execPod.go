package ws

import (
	client_go "Crd-End/client-go"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// ssh流式处理器
type streamHandler struct {
	wsConn *WsConnection
	resizeEvent chan remotecommand.TerminalSize
}

// web终端发来的包
type xtermMessage struct {
	MsgType string `json:"type"`	// 类型:resize客户端调整终端, input客户端输入
	Input string `json:"input"`	// msgtype=input情况下使用
	Rows uint16 `json:"rows"`	// msgtype=resize情况下使用
	Cols uint16 `json:"cols"`// msgtype=resize情况下使用
}

// executor回调获取web是否resize
func (handler *streamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <- handler.resizeEvent
	size = &ret
	return
}

// executor回调读取web端的输入
func (handler *streamHandler) Read(p []byte) (size int, err error) {
	var (
		msg *WsMessage
		xtermMsg xtermMessage
	)

	// 读web发来的输入
	if msg, err = handler.wsConn.WsRead(); err != nil {
		return
	}

	// 解析客户端请求
	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		return
	}

	//web ssh调整了终端大小
	if xtermMsg.MsgType == "resize" {
		// 放到channel里，等remotecommand executor调用我们的Next取走
		handler.resizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	} else if xtermMsg.MsgType == "input" {	// web ssh终端输入了字符
		// copy到p数组中
		size = len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return
}

// executor回调向web端输出
func (handler *streamHandler) Write(p []byte) (size int, err error) {
	var (
		copyData []byte
	)

	// 产生副本
	copyData = make([]byte, len(p))
	copy(copyData, p)
	size = len(p)
	err = handler.wsConn.WsWrite(websocket.TextMessage, copyData)
	return
}

//resp http.ResponseWriter, req *http.Request
func WsHandler(c *gin.Context) {
	var (
		wsConn *WsConnection
		restConf *rest.Config
		sshReq *rest.Request
		podName string
		podNs string
		containerName string
		executor remotecommand.Executor
		handler *streamHandler
		err error
	)
	//resp := c.Writer
	//req := c.Request

	podName = c.Request.FormValue("podname")
	containerName = c.Request.FormValue("containername")
	podNs = c.Request.FormValue("podns")

	// 得到websocket长连接
	wsConn, err = InitWebsocket(c);
	if err != nil {
		return
	}

	//// 获取k8s rest client配置 [dev]
	//if restConf, err = GetRestConf(); err != nil {
	//	goto END
	//}

	// 获取k8s rest client配置 [proc]
	if 	restConf, err = rest.InClusterConfig() ; err != nil {
		goto END
	}

	// URL长相:
	// https://172.18.11.25:6443/api/v1/namespaces/default/pods/nginx-deployment-5cbd8757f-d5qvx/exec?command=sh&container=nginx&stderr=true&stdin=true&stdout=true&tty=true
	sshReq = client_go.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(podNs).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   []string{"bash"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	// 创建到容器的连接
	if executor, err = remotecommand.NewSPDYExecutor(restConf, "POST", sshReq.URL()); err != nil {
		goto END
	}

	// 配置与容器之间的数据流处理回调
	handler = &streamHandler{ wsConn: wsConn, resizeEvent: make(chan remotecommand.TerminalSize)}
	if err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		goto END
	}
	return

END:
	fmt.Println(err)
	wsConn.WsClose()

}

// 获取k8s restful client配置
func GetRestConf() (restConf *rest.Config, err error) {
	var (
		kubeconfig []byte
	)

	// 读kubeconfig文件
	if kubeconfig, err = ioutil.ReadFile("config/admin.bak"); err != nil {
		goto END
	}

	// 生成rest client配置
	if restConf, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig); err != nil {
		goto END
	}
END:
	return
}