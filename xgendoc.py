#!/usr/bin/python

# Copyright 2015 The Goga Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import subprocess

def Cmd(command, verbose=False, debug=False):
    if debug:
        print '=================================================='
        print cmd
        print '=================================================='
    spr = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out = spr.stdout.read()
    err = spr.stderr.read().strip()
    if verbose:
        print out
        print err
    return out, err

odir  = 'doc/'
idxfn = odir+'index.html'
licen = open('LICENSE', 'r').read()

def header(title):
    return """<html>
<head>
<meta http-equiv=\\"Content-Type\\" content=\\"text/html; charset=utf-8\\">
<title>%s</title>
<link type=\\"text/css\\" rel=\\"stylesheet\\" href=\\"static/style.css\\">
<script type=\\"text/javascript\\" src=\\"static/godocs.js\\"></script>
<style type=\\"text/css\\"></style>
</head>
<body>
<div id=\\"page\\" class=\\wide\\">
<div class=\\"container\\">
""" % title

def footer():
    return """
<div id=\\"footer\\">
<br /><br />
<hr>
<pre class=\\"copyright\\">
%s</pre><!-- copyright -->
</div><!-- footer -->

</div><!-- container -->
</div><!-- page -->
</body>
</html>""" % licen

Cmd('echo "'+header('Goga &ndash; Documentation')+'" > '+idxfn)
Cmd('echo "<h1>Goga &ndash; Documentation</h1>" >> '+idxfn)
Cmd('echo "<h2 id=\\"pkg-index\\">Index</h2>\n<div id=\\"manual-nav\\">\n<dl>" >> '+idxfn)

Cmd('godoc -html github.com/cpmech/goga >> '+idxfn)

# fix links
Cmd("sed -i -e 's@/src/target@https://github.com/cpmech/goga/blob/master/@g' "+idxfn+"")
Cmd("sed -i -e 's@/src/github.com/cpmech/goga/@https://github.com/cpmech/goga/blob/master/@g' "+idxfn+"")

# fix links to subdirectories (harder to automate)
subdirs = ["data", "doc", "examples"]

for subdir in subdirs:
    Cmd("sed -i -e 's@<a href=\""+subdir+"/\">@<a href=\"https://github.com/cpmech/goga/tree/master/"+subdir+"\">@g' "+idxfn)

Cmd('echo "</dl>\n</div><!-- manual-nav -->" >> '+idxfn)
Cmd('echo "'+footer()+'" >> '+idxfn)
