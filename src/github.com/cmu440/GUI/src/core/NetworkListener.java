package core;

import java.io.IOException;
import java.io.InputStream;
import java.io.ObjectInputStream;
import java.net.Socket;
import java.net.UnknownHostException;

public class NetworkListener implements Runnable {
	private Socket s;
	private InputStream in;
	private Jeopardy j;
	public NetworkListener(Socket s,Jeopardy j){
		this.s=s;
		this.j=j;
		try {
			in =s.getInputStream();
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	private String jsonString(String command,String json){

				
		String withoutCommand= (json.replace(command, ""));
		return withoutCommand.replace("}\n", "");
	}
	@Override
	public void run() {
				byte[] buf=new byte[100];

				while (s.isConnected()){ 
					try {
						 in.read(buf);
						 
						 String msg=new String(buf,"UTF-8");
						 if(msg.startsWith("{\"Join\":")){
							 j.joined(jsonString("{\"Join\":",msg));
						 } else if (msg.startsWith("{\"Buzz")){
							 
						 } else if (msg.startsWith("{\"Answer")){
							 
						 } else if (msg.startsWith("{\"Join")){
							 
						 } else if (msg.startsWith("{\"Question")){
							 
						 } 
						 
					} catch (IOException e) {
						// TODO Auto-generated catch block
						e.printStackTrace();
					} 
			       
				}
		
	}

}
