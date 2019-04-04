# H

See:

~~~ txt
cluster.local:53 {
    kubernetes cluster.local {
        upstream
    }
}
example.local {
    forward . 10.100.0.10:53
}

. {
    forward . 8.8.8.8:53
}
~~~
