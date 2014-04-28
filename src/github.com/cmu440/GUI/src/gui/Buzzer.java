package gui;

import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;
import java.io.IOException;

import core.Jeopardy;

public class Buzzer implements ActionListener {
	private Jeopardy game;
	public Buzzer(Jeopardy game){
		this.game=game;
	}

	@Override
	public void actionPerformed(ActionEvent arg0) {
		try {
			game.buzz();
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}

	}

}
