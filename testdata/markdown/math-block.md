In conclusion, the following is true in Go:

$$
\textsf{ch := make(chan type, value)}
\left\{
\begin{array}{ll}
value == 0 & \rightarrow \textsf{unbuffered} \\
value >  0 & \rightarrow \textsf{buffer }{} value{} \textsf{ elements}
\end{array}
\right.
$$

More text.
