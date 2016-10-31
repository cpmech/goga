lcmd=latex -interaction=batchmode -halt-on-error
bcmd=bibtex
#pspdf=ps2pdf -dMaxSubsetPct=100 -dCompatibilityLevel=1.4 -dSubsetFonts=true -dEmbedAllFonts=true -dAutoFilterColorImages=false -dAutoFilterGrayImages=false -dColorImageFilter=/FlateEncode -dGrayImageFilter=/FlateEncode -dMonoImageFilter=/FlateEncode 
pspdf=ps2pdf

all: main

main:
	@echo "[1;32m################################### main #######################################[0m"
	$(lcmd)  main
	bibtex   main
	$(lcmd)  main
	dvips    main.dvi
	$(pspdf) main.ps

high:
	@echo "[1;32m################################### highlights #######################################[0m"
	$(lcmd)  highlights
	$(lcmd)  highlights
	dvips    highlights.dvi
	$(pspdf) highlights.ps
	rm       highlights.dvi
	rm       highlights.ps

cov:
	@echo "[1;32m################################### coverletter #######################################[0m"
	pdflatex coverletter
	pdflatex coverletter

coverletterR1:
	@echo "[1;32m################################### coverletterR1 #######################################[0m"
	pdflatex coverletterR1
	pdflatex coverletterR1

rev:
	@echo "[1;32m################################### revcom #######################################[0m"
	pdflatex revcom
	pdflatex revcom

cl:
	rm *.aux *.log *.bbl *.blg *.out *.toc *.tit *.spl *.ps *.dvi
