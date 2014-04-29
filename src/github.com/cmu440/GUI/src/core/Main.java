package core;

import java.io.FileNotFoundException;
import java.util.Random;

import gui.JeopardyPanel;

import javax.swing.JFrame;
import javax.swing.SwingUtilities;

public class Main {
	  private static final String NAME = "Jeopardy! 15440 Review";

	    public static void main(String[] args) {
	        SwingUtilities.invokeLater(new Runnable() {
	            @Override
	            public void run() {
	                try {
						createAndShowJeopardy();
					} catch (FileNotFoundException e) {
						// TODO Auto-generated catch block
						e.printStackTrace();
					}
	            }
	        });
	    }
	    private static void createAndShowJeopardy() throws FileNotFoundException {
	        // Create and set-up the window.
	        JFrame frame = new JFrame(NAME);
	        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
	        Random r=new Random();
	        int random=10000+r.nextInt(10000);
	        Jeopardy game = new Jeopardy("localhost:"+random,random);
	        
	        // Create and set up the content pane
	        JeopardyPanel gamePanel = new JeopardyPanel(game);
	        gamePanel.setOpaque(true);
	        frame.setContentPane(gamePanel);
	        frame.pack();
	        frame.setVisible(true);
	    }

}
