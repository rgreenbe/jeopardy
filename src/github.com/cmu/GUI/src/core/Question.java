package core;

import java.util.ArrayList;

public class Question {
	private String question;
	private ArrayList<String> options;
	private int answer;
	private int value;
	
	public Question(String question) {
		this.question=question;
		this.options=new ArrayList<String>();
	}
	public void setAnswer(int a){
		answer=a;
	}
	public void setOptions(String option){
		options.add(option);
	}
	public int answer(){
		return answer;
	}
	public String question(){
		return question;
	}
	public ArrayList<String> options(){
		return options;
	}
	public void setValue(int v) {
		value=v;
		
	}
	public int value() {
		return value;
	}

}
