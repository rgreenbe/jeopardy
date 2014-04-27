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
	private final String join,buzz,answer,question;
	public NetworkListener(Socket s,Jeopardy j){
		this.s=s;
		this.j=j;
		this.join="{\"Join\":";
		this.buzz="{\"Buzz\":";
		this.answer="{\"Answer\":";
		this.question="{\"Question\":";
				
		try {
			in =s.getInputStream();
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	private String jsonString(String command,String json){

				
		String withoutCommand= (json.replace(command, ""));
		withoutCommand=withoutCommand.replace("}\n", "");
		return withoutCommand.trim().replaceAll("\0", "");
	}
	@Override
	public void run() {
		byte[] buf=new byte[100];
	
		while (s.isConnected()){ 
			try {
				
				 in.read(buf);
				 
				 String msg=new String(buf,"UTF-8");
				 if(msg.startsWith(join)){
					 j.joined(jsonString(join,msg));
				 } else if (msg.startsWith(buzz)){
					 j.buzzed(jsonString(buzz,msg));
				 } else if (msg.startsWith(answer)){
					 j.answered(jsonString(answer,msg));
				 } else if (msg.startsWith(question)){
					 j.selectedQuestion(jsonString(question,msg));			 
				 } 
				 
			} catch (IOException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			} 
	       
		}

	}

}
