package gui;

import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.io.IOException;

import core.Jeopardy;

public class ChooseQuestion implements ActionListener {
	private int row, col;
	private Jeopardy game;

	public ChooseQuestion(int row, int col, Jeopardy game) {
		this.row = row;
		this.col = col;
		this.game = game;
	}

	@Override
	public void actionPerformed(ActionEvent arg0) {
		try {
			game.chooseQuestion(row, col);
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}

	}

}
