package core;


import java.util.Map;

import com.google.gson.Gson;


public class Jeopardy {
	private Map<String,Integer> players;
	private String name;
	private int gameID;
	private Gson gson;
	private GameInfo gameInfo;
	public Jeopardy(String json){
		gson=new Gson();
		this.gameInfo=gson.fromJson(json,GameInfo.class);
	}
	public GameInfo Info(){
		return gameInfo;
	}
	public void chooseQuestion(int row, int col) {
		System.out.println(Integer.toString(row)+Integer.toString(col));
		
	}
}
