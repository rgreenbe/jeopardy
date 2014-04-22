package gui;

import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;

import javax.swing.JButton;

import core.Jeopardy;

public class ChooseQuestion implements ActionListener {
	private int row,col;
	private JButton b;
	private Jeopardy game;
	public ChooseQuestion(int row,int col,JButton b,Jeopardy game){
		this.row=row;
		this.col=col;
		this.b=b;
		this.game=game;
	}

	@Override
	public void actionPerformed(ActionEvent arg0) {
		game.chooseQuestion(row,col);

	}

}
