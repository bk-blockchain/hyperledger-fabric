import java.io.File;
import java.util.Formatter;
import java.util.Scanner;

import java.io.FileNotFoundException;

public class read {
    public static void main(String[] args) {
        try {
            File x = new File("throughput.txt");
            Scanner sc = new Scanner(x);
            long time = sc.nextLong();
            long minTime=time,maxTime = time;
            long co=1;
            while (sc.hasNext()) {
                time = sc.nextLong();
                co+=1;
                if (minTime>time) minTime=time;
                if(maxTime<time) maxTime=time;
            }
            float th= (float)(co*1000)/ (float)(maxTime-minTime);
            System.out.println("throughput: "+th);
            sc.close();
        } catch (FileNotFoundException e) {
            System.out.println("Error");
        }

        try {
            File y = new File("latency.txt");
            Scanner scan = new Scanner(y);
            int time1=0 ;
            int sum =0, cout=0;
            while (scan.hasNext()) {
                time1=scan.nextInt();
                sum +=time1;
                cout+=1;
            }
            float la=(float)sum/(float)cout;
            System.out.println("latency : "+la);
            System.out.println("count : "+cout);
            scan.close();
        } catch (FileNotFoundException e) {
            System.out.println("Error");
        }
    }
}