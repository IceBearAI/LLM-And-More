/**
 * Copyright FunASR (https://github.com/alibaba-damo-academy/FunASR). All Rights
 * Reserved. MIT License  (https://opensource.org/licenses/MIT)
 */
/* 2021-2023 by zhaoming,mali aihealthx.com */

function WebSocketConnectMethod( config ) { //定义socket连接方法类

	
	var speechSokt;
	var connKeeperID;
	
	var msgHandle = config.msgHandle;
	var stateHandle = config.stateHandle;
			  
	this.wsStart = function () {
		let wsOrigin="";
		let token= JSON.parse(localStorage.getItem("user")).userInfo.token;
		var Uri = `${wsOrigin}/api/voice/asr/ws?X-Token=${token}`;
		// console.log("Uri="+Uri)
 
		if ( 'WebSocket' in window ) {
			speechSokt = new WebSocket( Uri ); // 定义socket连接对象
			speechSokt.onopen = function(e){onOpen(e);}; // 定义响应函数
			speechSokt.onclose = function(e){
			    console.log("onclose ws!");
			    //speechSokt.close();
				onClose(e);
				};
			speechSokt.onmessage = function(e){onMessage(e);};
			speechSokt.onerror = function(e){onError(e);};
			return 'broswerSupport';
		}
		else {
			return 'broswerNotSupport';
		}
	};
	
	// 定义停止与发送函数
	this.wsStop = function () {
		if(speechSokt != undefined) {
			console.log("stop ws!");
			speechSokt.close();
		}
	};
	
	this.wsSend = function ( oneData ) {
 
		if(speechSokt == undefined) return;
		if ( speechSokt.readyState === 1 ) { // 0:CONNECTING, 1:OPEN, 2:CLOSING, 3:CLOSED
 
			speechSokt.send( oneData );
 
			
		}
	};
	
	// SOCEKT连接中的消息与状态响应
	function onOpen( e ) {
		// 发送json
		var chunk_size = [5, 10, 5];
		var request = {
			chunk_size: chunk_size,
			"wav_name":  "h5",
			"is_speaking":  true,
			"chunk_interval":10,
			"itn":false,
			"mode":"online",
			
		};
		speechSokt.send(JSON.stringify(request));
		stateHandle('ws-success');
	}
	
	function onClose( e ) {
		stateHandle('ws-close');
	}
	
	function onMessage( e ) {
		msgHandle( e );
	}
	
	function onError( e ) {
		console.log("连接异常",e);
		stateHandle('ws-error');
	}
}