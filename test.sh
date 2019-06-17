#!/bin/sh

expectfailure() {
    echo "Expecting failure from $@:"
    env "$@" env2star -prefix 'a,b'
    if [ $? != 1 ]; then
	echo Failed
	exit 1
    fi
    echo
}

expectsuccess() {
    prefix="$1"; shift
    output="$1"; shift
    echo "$output: Expecting success from $@:"
    env "$@" env2star -prefix "$prefix" -output "$output"
    if [ $? != 0 ]; then
	echo Failed
	exit 1
    fi
    echo
}

## Failures

# a is both a terminating value and a map
expectfailure 'a=1' 'a.b=2'
expectfailure 'a.b=1' 'a=2'
# a is both a terminating value and an array
expectfailure 'a=1' 'a[0]=2'
expectfailure 'a[0]=1' 'a=2'
# a is both an array and a map
expectfailure 'a[0]=1' 'a.b=2'
expectfailure 'a.b=1' 'a[0]=2'

# a[0] is both a terminating value and a map
expectfailure 'a[0]=1' 'a[0].b=2'
expectfailure 'a[0].b=1' 'a[0]=2'
# a[0] is both a terminating value and an array
expectfailure 'a[0]=1' 'a[0][0]=2'
expectfailure 'a[0][0]=1' 'a[0]=2'
# a[0] is both an array and a map
expectfailure 'a[0][0]=1' 'a[0].b=2'
expectfailure 'a[0].b=1' 'a[0][0]=2'

# a[0][0] is both a terminating value and a map
expectfailure 'a[0][0]=1' 'a[0][0].b=2'
expectfailure 'a[0][0].b=1' 'a[0][0]=2'
# a[0][0] is both a terminating value and an array
expectfailure 'a[0][0]=1' 'a[0][0][0]=2'
expectfailure 'a[0][0][0]=1' 'a[0][0]=2'
# a[0][0] is both an array and a map
expectfailure 'a[0][0][0]=1' 'a[0][0].b=2'
expectfailure 'a[0][0].b=1' 'a[0][0][0]=2'

# a[0][0].b.c[0][0] is both a terminating value and a map
expectfailure 'a[0][0].b.c[0][0]=1' 'a[0][0].b.c[0][0].b=2'
expectfailure 'a[0][0].b.c[0][0].b=1' 'a[0][0].b.c[0][0]=2'
# a[0][0].b.c[0][0] is both a terminating value and an array
expectfailure 'a[0][0].b.c[0][0]=1' 'a[0][0].b.c[0][0][0]=2'
expectfailure 'a[0][0].b.c[0][0][0]=1' 'a[0][0].b.c[0][0]=2'
# a[0][0].b.c[0][0] is both an array and a map
expectfailure 'a[0][0].b.c[0][0][0]=1' 'a[0][0].b.c[0][0].b=2'
expectfailure 'a[0][0].b.c[0][0].b=1' 'a[0][0].b.c[0][0][0]=2'

# a[0][0].b.c[0][0].x.y.z is both a terminating value and a map
expectfailure 'a[0][0].b.c[0][0].x.y.z=1' 'a[0][0].b.c[0][0].x.y.z.b=2'
expectfailure 'a[0][0].b.c[0][0].x.y.z.b=1' 'a[0][0].b.c[0][0].x.y.z=2'
# a[0][0].b.c[0][0].x.y.z is both a terminating value and an array
expectfailure 'a[0][0].b.c[0][0].x.y.z=1' 'a[0][0].b.c[0][0].x.y.z[0]=2'
expectfailure 'a[0][0].b.c[0][0].x.y.z[0]=1' 'a[0][0].b.c[0][0].x.y.z=2'
# a[0][0].b.c[0][0].x.y.z is both an array and a map
expectfailure 'a[0][0].b.c[0][0].x.y.z[0]=1' 'a[0][0].b.c[0][0].x.y.z.b=2'
expectfailure 'a[0][0].b.c[0][0].x.y.z.b=1' 'a[0][0].b.c[0][0].x.y.z[0]=2'

# miscellaneous

# non number array indices
expectfailure 'a[0][a][0]=1'
expectfailure 'a[0]a=1'

## Successes

# empty prefix; outputs current environment
expectsuccess '' json

# random junk
expectsuccess a,b,c json 'a={}' 'b.d=[]' 'b.e=1.2' 'c[0][2].d=hey' 'c[0][1].d=false' 'c[1].abc=true' 'c[2]=null' 'd=notaprefix'
expectsuccess a,b,c yaml 'a={}' 'b.d=[]' 'b.e=1.2' 'c[0][2].d=hey' 'c[0][1].d=false' 'c[1].abc=true' 'c[2]=null' 'd=notaprefix'

expectsuccess a,d json 'a.b[0].c=2' 'd[0].b[0].c=2' 'd[0].b[1].d=3' 'd[0].b[2].x[0][0]=3' 'd[0].b[2].x[1][0].b=2'
expectsuccess a,d toml 'a.b[0].c=2' 'd[0].b[0].c=2' 'd[0].b[1].d=3' 'd[0].b[2].x[0][0]=3' 'd[0].b[2].x[1][0].b=2'
#expectsuccess a,d toml 'a.b[0]=1' 'a.b[1]=2'
echo hi


# doubly nested array whose elements are maps
expectsuccess a,b,c json 'a={}' 'b.d=[]' 'b.e=1.2' 'c[0][0].d=hey' 'c[0][1].d=false' 'c[1][0].abc=true' 'd=notaprefix'
expectsuccess a,b,c toml 'a={}' 'b.d=[]' 'b.e=1.2' 'c[0][0].d=hey' 'c[0][1].d=false' 'c[1][0].abc=true' 'd=notaprefix'

