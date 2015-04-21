

# Interpolation

## Interpolation in selector

$name: 'foo'


With suffix

    .#{ name }-suffix { }

Which should be

    .foo-suffix { }

With prefix and suffix:

    .prefix-#{ name }-suffix { }

Which will be compiled to:

    .prefix-foo-suffix {  }

    prefix#{ name }:hover { }
    prefix#{ name } { }

    pr#{...}efix#{ name } { }

Unknown selector:

    #{ foo } {

    }

    // property will end of ';' or '}'
    // selectors will end of '{', ','



## Interpolation in property name

    #{ prop_name + suffix }: 
    width-#{ prop_name + suffix }: 
    width-#{ prop_name + suffix }: 



## Solution

When lexing selectors, we just skip the interpolation in the selector. treat them as a part of selector and mark the token "ContainsInterpolation".

and we will parse it again in the parser to expand....


    [String] + [Interpolation] + [String]
                     |
                     |
                [Expression]
                     |
                   .....

