simpleParser
============

Using golang to implement a simple parser and generator for studying Compilers.

<p>
The script is somewhat like the JavaScript arithmetic statements.
Some of the simple statements like:
</p>
<pre>
var a;
a = 23 * 2;
var b;
b = a + 6 / 2;
echo b;
</pre>
For now the detailed description is:
<pre>
STMT :: VARI = EXPR ; | "var" VARI ; | "echo" EXPR ;
VARI :: [a-zA-Z][a-zA-Z0-9]*
EXPR :: EXPR + EXPR | EXPR - EXPR
EXPR :: EXPR * EXPR | EXPR / EXPR | EXPR % EXPR
EXPR :: ( EXPR ) | CONS | VARI
CONS :: [1-9][0-9]*(.[0-9]+)?
</pre>

There is also a state chart of it(Which is not so-optimized):
<pre>
       "var"      VARI        ;
      /----->(10)----->(11)------>(12)
     /                         ; ^
    /       "echo"              /
 (0)-------------------->(100)-/
   \                  /   /   ^
    \ VARI        =  /   /+-*/%\
     \---->(200)----/    \-----/
</pre>
