package gui;

import java.awt.BorderLayout;
import java.awt.Color;
import java.awt.Dimension;
import java.awt.GridLayout;
import java.util.List;

import javax.swing.JLabel;
import javax.swing.JPanel;

import core.GameInfo;
import core.Jeopardy;

public class JeopardyPanel extends JPanel {
	private final Jeopardy j;
    private JPanel headPanel,footPanel,gridPanel;
	  public JeopardyPanel(Jeopardy game) {

		 j=game;
		 headPanel=createHeadPanel();
		 gridPanel=createGrid();
		 footPanel=createFootPanel();
		 initGui();
	  }

	private JPanel createGrid() {
		
		GameInfo gameInfo=j.Info();
		JPanel grid=new JPanel();

		List<List<Boolean>> board=gameInfo.board();
		int rows=board.size();
		int cols=board.get(0).size();
		setLayout(new GridLayout(rows,cols));
		for (int row=0;row<rows;row++){
			for(int col=0;col<cols;col++){
				JLabel p=new JLabel("HEY");
				if( board.get(row).get(col)){
					p.setOpaque(true);
					p.setBackground(Color.black);
				}
				grid.add(p);
			}
		}
		return grid;
	}

	private void initGui() {
		setLayout(new BorderLayout());
        Color darkCyan = new Color(10, 135, 171);
        setPreferredSize(new Dimension(800,700));
        add(headPanel,BorderLayout.NORTH);
        add(gridPanel,BorderLayout.CENTER);
        add(footPanel,BorderLayout.SOUTH);
        setBackground(darkCyan);
        setOpaque(true);

			
	}
    private JPanel createFootPanel() {
        JPanel panel = new JPanel();
        JLabel copyright = new JLabel("Copyright © 2014 Gode Greenberg 15-440 Jeopardy!");
        panel.add(copyright);
        return panel;
    }
    private JPanel createHeadPanel() {
        JPanel panel = new JPanel();
        JLabel copyright = new JLabel("Welcome To Jeopardy! Review for the 15440 Final");
        panel.add(copyright);
        return panel;
    }

}
