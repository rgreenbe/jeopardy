package core;


import java.io.IOException;
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
	public Jeopardy(){
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
	public void joinGame(){
		System.out.println("Join:");
		Scanner sc = new Scanner(System.in);
		System.out.println("waiting for response");
		
		try {
			System.in.read();
			System.out.println("GOT INPUT ");
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		//You need 
	}
	
}
