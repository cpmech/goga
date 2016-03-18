# Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from DrawMsh import *
from gosl    import *

SetForEps(1.5, 350)
d = DrawMsh('frame2d.msh')
d.txt_fs = 8
d.linEC = {-1:'#9b0000', -2:'#55b34e', -3:'#1800a9', -4:'#84c1ff', -5:'#ff7800'}
d.linLW = 3
d.ftagC = 'black'
d.cidC = 'black'
d.vidC = 'black'
d.cid_txt_dx = {-1:0.2, -2:0.2, -3:0.2, -4:0.0, -5:0.0}
d.cid_txt_dy = {-1:0.0, -2:0.0, -3:0.0, -4:0.5, -5:0.5}
d.draw(show_vtags=0, show_ctags=1, tags_only=1, show_vids=0, show_cids=0,
        show_points=0, show_ftags=0, noGll=0, use_ec_txt=1, pos_tags=1)
text(   0,-1.2,   '$0$',fontsize=7,ha='center')
text( 7.5,-1.2, '$7.5$',fontsize=7,ha='center')
text(  11,-1.2,  '$11$',fontsize=7,ha='center')
text(18.5,-1.2,'$18.5$',fontsize=7,ha='center')
axvline(   0,color='grey',ls='--',lw=0.8)
axvline( 7.5,color='grey',ls='--',lw=0.8)
axvline(  11,color='grey',ls='--',lw=0.8)
axvline(18.5,color='grey',ls='--',lw=0.8)
SetXnticks(0)
SetYnticks(13)
for y in linspace(4,48,12):
    Arrow(-5,y,-1.8,y, st='->', zorder=10)
    text(-4.4,y+0.2,'$P$')
Circle(18.5, 48, 0.7)
text(19.0, 48.8, 'A', size=12)
axis([-5,20,axis()[2],axis()[3]])
gca().xaxis.labelpad = 15
Save("frame2d.eps")
