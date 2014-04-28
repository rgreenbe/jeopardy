package gui;

import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.io.IOException;

import core.Jeopardy;

public class ChooseAnswer implements ActionListener {
	private int choice;
	private Jeopardy game;
	public ChooseAnswer(int choice,Jeopardy game){
		this.choice=choice;
		this.game=game;
	}
	@Override
	public void actionPerformed(ActionEvent arg0) {
		try {
			game.chooseAnswer(choice);
		} catch (IOException e) {
			e.printStackTrace();
		}

	}

}
