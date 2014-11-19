reset
set terminal epslatex size 8.89cm,6.65cm color colortext
set output "open.tex"
set ylabel ""
set grid
set xrange [0:1]
set yrange [0:1]
plot "open.dat" 
set output
!cp open.tex ~/shared/
