from DrawMsh  import *
from gosl     import *
from pylab    import xticks, yticks

fnk = 'ground10'

SetForEps(0.7,280)

d = DrawMsh('%s.msh' % fnk)
d.draw(show_points=0, show_ftags=0, show_ctags=0, noGll=True)
xticks([0, 60, 120, 180, 240, 300, 360, 420, 480, 540, 600, 660, 720])
yticks([0, 60, 120, 180, 240, 300, 360])
axis([-50,770,-100,360])
Arrow(360,0,360,-100, sc=12, clip_on=0, zorder=100)
Arrow(720,0,720,-100, sc=12, clip_on=0, zorder=100)
text(350,-50,'100',ha='right')
text(710,-50,'100',ha='right')
text(-40,370,'fully fixed')
text(-40,-10,'fully fixed',va='top')
#grid(color='gray')
Gll("$x$", "$y$", "")
Save('mesh-%s.eps' % fnk)
