package core;

import java.util.List;

public class GameInfo {
	private int playerID;
	private int gameID;
	private List<Integer> players;
	private List<List<Integer>> board;
	private List<String> categories;
	private Boolean canControl;
	public GameInfo(){}
	public int playerID(){
		return playerID;
	}
	public int gameID(){
		return gameID;
	}
	public List<Integer> players(){
		return players;
	}
	public List<List<Integer>> board(){
		return board;
	}
	public List<String> Categories() {
		return categories;
		
	}
	
}
