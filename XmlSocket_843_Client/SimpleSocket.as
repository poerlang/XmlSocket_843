package
{
	import flash.events.Event;
	import flash.events.IOErrorEvent;
	import flash.events.ProgressEvent;
	import flash.events.SecurityErrorEvent;
	import flash.net.Socket;
	import flash.utils.ByteArray;
	
	/**
	 * 向服务端发送安全策略请求
	 */
	/**
	 * 
	 * @author Administrator
	 */
	/**
	 * 
	 * @author Administrator
	 */
	public class SimpleSocket extends Socket
	{
		/**
		 *构造函数
		 */
		public function SimpleSocket()
		{
			start();
		}

		/**
		 * 启动自定义socket
		 */
		public function start():void
		{
			configureListeners();
			super.connect("127.0.0.1", 843);
		}
		
		/**
		 * 发送安全策略请求
		 */
		public function sendMessage():void
		{
			var b:ByteArray = new ByteArray();
			b.writeUTFBytes("<policy-file-request/>\0");
			this.writeBytes(b);
			this.flush();
		}
		/**
		 * 配置socket监听事件
		 */
		private function configureListeners():void
		{
			addEventListener(Event.CLOSE, closeHandler);
			addEventListener(Event.CONNECT, connectHandler);
			addEventListener(IOErrorEvent.IO_ERROR, ioErrorHandler);
			addEventListener(SecurityErrorEvent.SECURITY_ERROR, securityErrorHandler);
			addEventListener(ProgressEvent.SOCKET_DATA, socketDataHandler);
		}
		
		override public function close():void
		{
			super.close();
			trace("self close");
		}
		
		/**
		 * 服务端关闭后触发
		 */
		private function closeHandler(event:Event):void
		{
			trace("conn close");
		}
		
		/**
		 * 侦测到链接已贯通
		 */
		private function connectHandler(event:Event):void
		{
			sendMessage();
		}

		/**
		 * IO异常
		 */
		private function ioErrorHandler(event:IOErrorEvent):void
		{
			new Error(event.toString());
			try{
				this.close();
			}catch (e:Error){
				//
			}
		}

		/**
		 * 安全异常
		 */
		private function securityErrorHandler(event:SecurityErrorEvent):void
		{
			trace(event.toString());
		}

		/**
		 * 收到服务端反馈
		 */
		private function socketDataHandler(event:ProgressEvent):void
		{
			var bytes:ByteArray=new ByteArray();
			this.readBytes(bytes, 0, this.bytesAvailable);
			trace("收到服务器返回：");
			trace(bytes.readUTFBytes(bytes.length));
		}
	}
}
