package core;

import java.io.File;
import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.Scanner;

public class Questions {
	private ArrayList<String> q;
	public Questions(){
		q=new ArrayList<String>();
		String file=("assets/words.txt");
		Scanner s;
		String question;
		int i =0;
		try {
			s = new Scanner (new File(file));
		while(s.hasNextLine()){
			question=s.nextLine();
			q.add(question  );
			
			}
		} catch (FileNotFoundException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
}

