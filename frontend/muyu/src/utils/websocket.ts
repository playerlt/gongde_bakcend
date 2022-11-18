/**
 * 定义回调函数
 * @author ERRUI
 */
 class SocketService {
    static instance: any = null;
  
    /**
     * 单例模式
     */
    static get Instance() {
      if (SocketService.instance) {
        return SocketService.instance;
      }
      return (SocketService.instance = new SocketService());
    }
    // 和服务端连接的socket对象
    ws: any = null;
  
    // 存储回调函数
    callBackMapping = {};
  
    // 标识是否连接成功
    connected = false;
  
    // 记录重试的次数
    sendRetryCount = 0;
  
    // 重新连接尝试的次数
    connectRetryCount = 0;
    /**
     * 心跳检测机制:20秒发一个检测消息过去
     */
    sendFixHeart() {
      let sendFixHeartTimer: any = null;
      clearInterval(sendFixHeartTimer);
      sendFixHeartTimer = setInterval(() => {
        this.ws.send(JSON.stringify({}));
      }, 20000);
    }
    /**
     * websocket连接
     */
    connect() {
      // 连接服务器
      if (!window.WebSocket) {
        return console.log('您的浏览器不支持WebSocket');
      }
      this.ws = new WebSocket('ws://127.0.0.1:20081/getCount');
      // 连接成功的事件
      this.ws.onopen = () => {
        console.log('连接服务端成功了,ws状态是' + this.ws.readyState);
        this.connected = true;
        // 重置重新连接的次数
        this.connectRetryCount = 0;
        this.sendFixHeart();
      };
  
      // 当连接成功之后, 服务器关闭的情况 ｜  连接服务端失败
      this.ws.onclose = () => {
        if (this.connectRetryCount >= 10) {
          this.ws.close();
          console.log('重新连接服务器失败');
        } else {
          this.connected = false;
          this.connectRetryCount++;
          setTimeout(() => {
            this.close();
            this.connect();
            console.log('连接服务端失败第' + this.connectRetryCount + '次，1秒后重试');
          }, 1000);
        }
      };
  
      // websocket连接错误
      this.ws.onerror = function () {
        console.log('Error，websocket连接错误');
      };
    }
    send(data: any) {
      // 判断此时此刻有没有连接成功
      if (this.connected) {
        this.sendRetryCount = 0;
  
        this.ws.send(JSON.stringify(data));
      } else {
        this.sendRetryCount++;
        setTimeout(() => {
          this.send(data);
          console.log('已经重新发送第:' + this.sendRetryCount + '次');
        }, 1000);
      }
    }
  
    /**
     * 关闭websocket连接
     */
    close() {
      this.ws.close();
    }
  }
  export default SocketService;  