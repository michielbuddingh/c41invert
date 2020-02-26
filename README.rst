c41invert
=========

c41invert is a command-line tool to quickly convert scans of
orange-backed colour negatives into positives.  For me its main
selling point is its lack of knobs to tweak; it uses sensible defaults
to get sensible results; if a certain picture deserves a more
perfectionist approach, there is a lot of graphics software that will
help you achieve that.  This tool is meant to give you extra time
to use them.

It uses a similar technique to negfix8_, although I wasn't aware of
its existence at the time.  

the Approach
~~~~~~~~~~~~

The tool samples the central section of the image, creating a
histogram of colours for each colour channel.  It then picks a
suitably 'dark' and 'light' colour (the first and ninetynineth
percentile, respectively)  

How to use
~~~~~~~~~~

c41invert convert -i inputfile.name -o outputfile.jpg

Output files are always written as JPEG at quality setting 95.

The option -s-curve uses a sigmoid function (an S-shaped curve) rather
than a linear function; you might like it, I don't.

.. _negfix8: https://sites.google.com/site/negfix/howto
