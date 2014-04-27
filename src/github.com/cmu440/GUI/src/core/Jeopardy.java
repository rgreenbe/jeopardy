package core;


import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Map;
import java.util.Scanner;

import com.google.gson.Gson;


public class Jeopardy {
	private Map<String,Integer> players;
	private int playerID;
	private int gameID;
	private Gson gson;
	private GameInfo gameInfo;
	private int turn;
	private String hostport;
	private Socket s;
	private OutputStream out; 
	private ArrayList<Question> questions;
	private GameChangeListener listener;
	private final int rows,cols;
	public Jeopardy(String hostport) throws FileNotFoundException{
		this.hostport=hostport;
		this.gameID=0;
		this.playerID=0;
		rows=4;
		cols=5;
		this.questions=(new Questions()).questions();
		try {
			s = new Socket("localhost",8080);
			out = s.getOutputStream();
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
	private String makeJson(String command,String jsonObject){
		return ("{\""+command+"\":"+jsonObject+"}\n");
			
	}
	public GameInfo Info(){
		return gameInfo;
	}
	public void chooseQuestion(int row, int col) throws IOException {
		QuestionArgs qa=new QuestionArgs(playerID,gameID,row,col);
		System.out.println(gson.toJson(qa)+"SIZE: "+gson.toJson(qa).length());
		out.write(makeJson("Question",gson.toJson(qa)).getBytes());
	}
	public void buzz(){
		turn++;
	}
	public void answerQuestion(){
	}

	public void joinGame() throws IOException{
		JoinArgs join =new JoinArgs(hostport);
		System.out.println("joining game");
		out.write(makeJson("Join",gson.toJson(join)).getBytes());

	}
	public void addListener(GameChangeListener l){
		this.listener=l;
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
	private int qIndex(int row,int col){
		return (rows*row+col);
	}
	public synchronized void selectedQuestion(String json){
		System.out.println(json+ "SIZE: "+json.length());
		QuestionArgs q=gson.fromJson(json, QuestionArgs.class);
		System.out.println("row: "+Integer.toString(q.row())+"col: ");
		int r=q.row();
		int c=q.col();
		listener.selectQuestion(questions.get(qIndex(r,c)), r, c);
		
	}
	
}
