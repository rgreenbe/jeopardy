package core;

public class BuzzArgs {

	private int gameID;
	private int playerID;
	private int turn;
	public BuzzArgs(int gameID,int playerID,int turn){
		this.turn=turn;
		this.playerID=playerID;
		this.gameID=gameID;
	}
	public int playerID(){
		return playerID;
	}
	public int gameID(){
		return gameID;
	}
	public int turn(){
		return turn;
	}
}
