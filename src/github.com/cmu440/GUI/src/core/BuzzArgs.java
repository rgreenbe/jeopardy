package core;

public class BuzzArgs {

	private int gameID;
	private int playerID;
	private int turn;
	public BuzzArgs(int gameID,int playerID,int turn){
		this.turn=turn;
	}
	public int playerID(){
		return playerID;
	}
	
}
