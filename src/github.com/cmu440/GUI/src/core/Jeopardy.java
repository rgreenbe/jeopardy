package core;


import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.net.UnknownHostException;
import java.util.Map;
import java.util.Scanner;

import com.google.gson.Gson;


public class Jeopardy {
	private Map<String,Integer> players;
	private String name;
	private int gameID;
	private Gson gson;
	private GameInfo gameInfo;
	private int turn;
	private String hostport;
	private Socket s;

	public Jeopardy(String hostport){
		this.hostport=hostport;
		try {
			s = new Socket("localhost",8080);
			(new Thread(new NetworkListener(s,this))).start();
		} catch (UnknownHostException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		String json="{PlayerID:3,Test:[true,true],"
        		+ "board:[[200,200,200,200,200],"
        		+ "[400,400,400,400,400],"
        		+ "[600,600,600,600,600],"
        		+ "[800,800,800,800,800]],categories:[poop,hello,world,mgode,rgreenbe]"
        		+ "}";
		gson=new Gson();
		turn =0;
		this.gameInfo=gson.fromJson(json,GameInfo.class);
	}
	public GameInfo Info(){
		return gameInfo;
	}
	public void chooseQuestion(int row, int col) {
		System.out.println(Integer.toString(row)+Integer.toString(col));
		
	}
	public void buzz(){
		turn++;
	}
	public void answerQuestion(){
	}
	private String makeJson(String command,String jsonObject){
		return ("{\""+command+"\":"+jsonObject+"}\n");
			
	}
	public void joinGame() throws IOException{
		JoinArgs join =new JoinArgs(hostport);
		OutputStream out = s.getOutputStream();
		System.out.println("joining game");
		out.write(makeJson("Join",gson.toJson(join)).getBytes());

	}
	public synchronized void joined(String json){
		System.out.println(json);
		//gson.fromJson(json, JoinRep.class);
	}
	public synchronized void buzzed(String json){
		
	}
	public synchronized void answered(String json){
		
	}
	public synchronized void gameStarting(String json){
		
	}
	public synchronized void selectedQuestion(String json){
		
	}
	
}
