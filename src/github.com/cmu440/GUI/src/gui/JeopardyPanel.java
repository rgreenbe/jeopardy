package gui;

import java.awt.BorderLayout;
import java.awt.Color;
import java.awt.Dimension;
import java.awt.GridLayout;
import java.util.List;

import javax.swing.JButton;
import javax.swing.JLabel;
import javax.swing.JPanel;
import javax.swing.border.LineBorder;

import core.GameInfo;
import core.Jeopardy;

public class JeopardyPanel extends JPanel {
	private final Jeopardy j;
	private final int boardWidth;
	private final int boardLength;
	private  GameInfo gameInfo;
    private JPanel headPanel,footPanel,gridPanel;
	  public JeopardyPanel(Jeopardy game) {
		 boardWidth=800;
		 boardLength=700;
		 j=game;
		 gameInfo=game.Info();
		 headPanel=createHeadPanel();
		 gridPanel=createGrid();
		 footPanel=createFootPanel();

		 initGui();
	  }

	private JPanel createGrid() {
		JPanel grid=new JPanel();
		
		List<List<Integer>> board=gameInfo.board();
		int padding=5;
		int rows=board.size();
		int cols=board.get(0).size();
		grid.setMaximumSize(new Dimension(500,500));
		grid.setLayout(new GridLayout(rows,cols));
		for (int row=0;row<rows;row++){
			for(int col=0;col<cols;col++){
				int val= board.get(row).get(col);
				JPanel box=new JPanel();
				//.box.setPreferredSize(new Dimension((boardWidth/cols-padding,boardLength/rows-padding));
				JButton p=new JButton();
				if(val!=0){
					p.setOpaque(true);
					//p.setBorder(new LineBorder(Color.BLACK,2));

					grid.add(p);
					p.setMaximumSize(new Dimension((500)/rows-padding,(500)/cols-padding));
					p.setPreferredSize(new Dimension((500)/rows-padding,(500)/cols-padding));
					p.setText(Integer.toString(val));
				}
				//grid.add(box);

			}
		}
		return grid;
	}

	private JPanel addCategories(JPanel grid) {
		List<String> categories =gameInfo.Categories();
		grid.setLayout(new GridLayout(1,categories.size()));
		for(int i=0;i<categories.size();i++){
			JLabel category=new JLabel(categories.get(i),JLabel.CENTER);
			category.setMaximumSize(new Dimension(10,10));
			grid.add(category);
		}
		return grid;
		
	}

	private void initGui() {
		setLayout(new BorderLayout());
        Color darkCyan = new Color(10, 135, 171);
        setPreferredSize(new Dimension(boardWidth,boardLength));
        add(headPanel,BorderLayout.NORTH);
        add(gridPanel,BorderLayout.CENTER);
        add(footPanel,BorderLayout.SOUTH);
        setBackground(darkCyan);
        setOpaque(true);

			
	}
    private JPanel createFootPanel() {
        JPanel panel = new JPanel();
        JLabel copyright = new JLabel("Copyright © 2014 Gode Greenberg 15-440 Jeopardy!",JLabel.CENTER);
        panel.add(copyright);
        return panel;
    }
    private JPanel createHeadPanel() {
        JPanel panel = new JPanel();
        panel.setLayout(new BorderLayout());
        JPanel categoryGrid=new JPanel();
        JPanel categories= addCategories(categoryGrid);
        
        
        JLabel copyright = new JLabel("Welcome To Jeopardy! Review for the 15440 Final");
        panel.add(copyright,BorderLayout.NORTH);
        panel.add(categories,BorderLayout.SOUTH);
        return panel;
    }

}
