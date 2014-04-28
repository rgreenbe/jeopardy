package core;

public class QuestionArgs {
	private int playerID;
	private int gameID;
	private int row;
	private int col;
	public QuestionArgs(int playerID,int gameID,int row,int col){
		this.playerID=playerID;
		this.gameID=gameID;
		this.row=row;
		this.col=col;
	}
	public int playerID(){
		return playerID;
	}
	public int gameID(){
		return gameID;
	}
	public int row(){
		return row;
	}
	public int col(){
		return col;
	}
}
