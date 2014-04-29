package core;


import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.ServerSocket;
import java.net.Socket;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Map;
import java.util.Scanner;

import com.google.gson.Gson;


public class Jeopardy {
	private ArrayList<Integer> players;
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
	private Question currentQ;
	private final int rows,cols;
	public Jeopardy(String hostport) throws FileNotFoundException{
		this.hostport=hostport;
		this.gameID=0;
		this.playerID=0;
		rows=5;
		cols=4;
		players=new ArrayList<Integer>();
		players.add(0);
		this.questions=(new Questions()).questions();
		try {
			s = new Socket("localhost",8080);
			out = s.getOutputStream();
			ServerSocket readSocket = new ServerSocket(9090);
			(new Thread(new NetworkListener(readSocket,this))).start();
		} catch (UnknownHostException e) {
			e.printStackTrace();
		} catch (IOException e) {
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
		System.out.println("Choosing Question: "+Integer.toString(qIndex(row,col)));
		out.write(makeJson("Question",gson.toJson(qa)).getBytes());
	}
	public void buzz() throws IOException{
		turn++;
		BuzzArgs buzz=new BuzzArgs(gameID,playerID,turn);
		out.write(makeJson("Buzz",gson.toJson(buzz)).getBytes());
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
		System.out.println("Buzzed "+json);
		BuzzArgs b=gson.fromJson(json, BuzzArgs.class);
		System.out.println(json);
		listener.buzzedIn(b.playerID());
	}
	public synchronized void answered(String json) throws InterruptedException{
		AnswerArgs a=gson.fromJson(json, AnswerArgs.class);
		System.out.println("Question value: "+Integer.toString(a.scoreChange()));
		int score=players.get(a.playerID());
		score+=a.scoreChange();
		players.add(a.playerID(),score);
		listener.answeredQuestion(a.playerID(), score,a.choice());
		
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
		currentQ=questions.get(qIndex(r,c));
		listener.selectQuestion(currentQ, r, c);
		
	}
	public void chooseAnswer(int choice) throws IOException {
		int scoreChange=currentQ.value();
		if(currentQ.answer()!=choice){
			scoreChange*=-1;
		}
		AnswerArgs a=new AnswerArgs(gameID,playerID,scoreChange,choice);
		out.write(makeJson("Answer",gson.toJson(a)).getBytes());
		
	}
	public int playerID() {
		return playerID;
	}
	public void addPlayers(int id){
		players.add(id);
	}
	public ArrayList<Integer> players(){
		return players;
	}
	public Question currentQuestion(){
		return currentQ;
	}
	
}
