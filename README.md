[![Go Report Card](https://goreportcard.com/badge/github.com/MicahParks/keyfunc)](https://goreportcard.com/report/github.com/MicahParks/keyfunc) [![Go Reference](https://pkg.go.dev/badge/github.com/MicahParks/keyfunc.svg)](https://pkg.go.dev/github.com/MicahParks/keyfunc)

# keyfunc

Purpose of this package is to provide a
[`jwt.Keyfunc`](https://pkg.go.dev/github.com/dgrijalva/jwt-go@v3.2.0+incompatible#Keyfunc) for the
[github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) package and its popular forks using a JSON Web Key
Set (JWKs) for parsing and verifying JSON Web Tokens (JWTs).

It's common for an identity provider, such as [Keycloak](https://www.keycloak.org/)
or [Amazon Cognito (AWS)](https://aws.amazon.com/cognito/) to expose a JWKs via an HTTPS endpoint. This package has the
ability to consume that JWKs and produce a
[`jwt.Keyfunc`](https://pkg.go.dev/github.com/dgrijalva/jwt-go@v3.2.0+incompatible#Keyfunc). It is important that a JWKs
endpoint is using HTTPS to ensure the keys are from the correct trusted source.

This repository has the following dependencies:
* [github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)
* [github.com/form3tech-oss/jwt-go](https://github.com/form3tech-oss/jwt-go)

`jwt.Keyfunc` signatures are imported from these, implemented, then exported as methods.

## Supported Algorithms

It is recommended to only use this package for asymmetric signing keys. If you are using HMAC signing keys, this Go
package may be unnecessary as the algorithm is symmetric, meaning the key is pre-shared. In this case a JWKs is likely
not be the best solution.

Currently, this package supports JWTs signed with an `alg` that matches one of the following:
* `ES256`
* `ES384`
* `ES512`
* `PS256`
* `PS384`
* `PS512`
* `RS256`
* `RS384`
* `RS512`

Additionally, the supported elliptical curve types are below:
* `P-256`
* `P-384`
* `P-521`

If there are cryptographic algorithms, curve types, or something else already standardized that you'd like supported in
this Go package, please open an issue or pull request.

## Basic usage

For complete examples, please see the `examples` directory.

```go
import "github.com/MicahParks/keyfunc"
```

### Preconditions: Acquire the JWKs URL or JSON

A JWKs URL is not required, one can be created directly from JSON with the
[`keyfunc.New`](https://pkg.go.dev/github.com/MicahParks/keyfunc#New) function.

```go
// Get the JWKs URL from an environment variable.
jwksURL := os.Getenv("JWKS_URL")

// Confirm the environment variable is not empty.
if jwksURL == "" {
	log.Fatalln("JWKS_URL environment variable must be populated.")
}
```

### Step 1: Create the JWKs
Via HTTP:
```go
// Create the JWKs from the resource at the given URL.
jwks, err := keyfunc.Get(jwksURL)
if err != nil {
	log.Fatalf("Failed to get the JWKs from the given URL.\nError: %s", err.Error())
}
```
Via JSON:
```go
// Get the JWKs as JSON.
var jwksJSON = json.RawMessage(`{"keys":[{"kid":"zXew0UJ1h6Q4CCcd_9wxMzvcp5cEBifH0KWrCz2Kyxc","kty":"RSA","alg":"PS256","use":"sig","n":"wqS81x6fItPUdh1OWCT8p3AuLYgFlpmg61WXp6sp1pVijoyF29GOSaD9xE-vLtegX-5h0BnP7va0bwsOAPdh6SdeVslEifNGHCtID0xNFqHNWcXSt4eLfQKAPFUq0TsEO-8P1QHRq6yeG8JAFaxakkaagLFuV8Vd_21PGJFWhvJodJLhX_-Ym9L8XUpIPps_mQriMUOWDe-5DWjHnDtfV7mgaOxbBvVo3wj8V2Lmo5Li4HabT4MEzeJ6e9IdFo2kj_44Yy9osX-PMPtu8BQz_onPgf0wjrVWt349Rj6OkS8RxlNGYeuIxYZr0TOhP5F-yEPhSXDsKdVTwPf7zAAaKQ","e":"AQAB","x5c":["MIICmzCCAYMCBgF4HR7HNDANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjEwMzEwMTcwOTE5WhcNMzEwMzEwMTcxMDU5WjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDCpLzXHp8i09R2HU5YJPyncC4tiAWWmaDrVZenqynWlWKOjIXb0Y5JoP3ET68u16Bf7mHQGc/u9rRvCw4A92HpJ15WyUSJ80YcK0gPTE0Woc1ZxdK3h4t9AoA8VSrROwQ77w/VAdGrrJ4bwkAVrFqSRpqAsW5XxV3/bU8YkVaG8mh0kuFf/5ib0vxdSkg+mz+ZCuIxQ5YN77kNaMecO19XuaBo7FsG9WjfCPxXYuajkuLgdptPgwTN4np70h0WjaSP/jhjL2ixf48w+27wFDP+ic+B/TCOtVa3fj1GPo6RLxHGU0Zh64jFhmvRM6E/kX7IQ+FJcOwp1VPA9/vMABopAgMBAAEwDQYJKoZIhvcNAQELBQADggEBALILq1Z4oQNJZEUt24VZcvknsWtQtvPxl3JNcBQgDR5/IMgl5VndRZ9OT56KUqrR5xRsWiCvh5Lgv4fUEzAAo9ToiPLub1SKP063zWrvfgi3YZ19bty0iXFm7l2cpQ3ejFV7WpcdLJE0lapFdPLo6QaRdgNu/1p4vbYg7zSK1fQ0OY5b3ajhAx/bhWlrN685owRbO5/r4rUOa6oo9l4Qn7jUxKUx4rcoe7zUM7qrpOPqKvn0DBp3n1/+9pOZXCjIfZGvYwP5NhzBDCkRzaXcJHlOqWzMBzyovVrzVmUilBcj+EsTYJs0gVXKzduX5zO6YWhFs23lu7AijdkxTY65YM0="],"x5t":"IYIeevIT57t8ppUejM42Bqx6f3I","x5t#S256":"TuOrBy2NcTlFSWuZ8Kh8W8AjQagb4fnfP1SlKMO8-So"},{"kid":"ebJxnm9B3QDBljB5XJWEu72qx6BawDaMAhwz4aKPkQ0","kty":"EC","alg":"ES512","use":"sig","crv":"P-521","x":"YQ95Xj8MTzcHytbU1h8YkCN2kdEQA7ThuZ1ctB9Ekiw6tlM9RwL62eQvzEt4Rz8qN69uRqgU9RzxQOkSU5xVvyo","y":"SMMuP3QnAPHtx7Go2ARsG3NBaySWBLmVvS8s2Ss7Vm_ISWenNbdjKOsY1XvtiQz5scGzWDCEUoZzgV8Ve1mLOV0"},{"kid":"TVAAet63O3xy_KK6_bxVIu7Ra3_z1wlB543Fbwi5VaU","kty":"EC","alg":"ES384","use":"sig","crv":"P-384","x":"Pik2o5as-evijFABH5p6YLXHnWw8iQ_N1ummPY1c_UgG6NO0za-gNOhTz2-tsd_w","y":"e98VSff71k19SY_mHgp3707lgQVrhfVpiGa-sGaKxOWVpxd2jWMhB0Q4RpSRuCp5"},{"kid":"arlUxX4hh56rNO-XdIPhDT7bqBMqcBwNQuP_TnZJNGs","kty":"RSA","alg":"RS512","use":"sig","n":"hhtifu8LL3ICE3BAX5l1KZv6Lni0lhlhBusSfepnpxcb4C_z2U71cQTnLY27kt8WB4bNG6e5_KMx9K3xUdd3euj9MCq8vytwEPieeHE1KXQuhJfLv017lhpK_dRMOHyc-9-50YNdgs_8KWRkrzjjuYrCiO9Iu76n5319e-SC8OPvNUglqxp2N0Sp2ltne2ZrpN8T3OEEXT62TSGmLAVopRGw5gllNVrJfmEyZJCRrBM6s5CQcz8un0FjkAAC4DI6QD-eBL0qG3_NR0hQvR1he2o4BLwjOKH45Pk_jj-eArp-DD6Xq6ABQVb5SNOSdaxl5lnmuotRoY3G5d9YSl-K3w","e":"AQAB","x5c":["MIICmzCCAYMCBgF4HSCcDzANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjEwMzEwMTcxMTE5WhcNMzEwMzEwMTcxMjU5WjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCGG2J+7wsvcgITcEBfmXUpm/oueLSWGWEG6xJ96menFxvgL/PZTvVxBOctjbuS3xYHhs0bp7n8ozH0rfFR13d66P0wKry/K3AQ+J54cTUpdC6El8u/TXuWGkr91Ew4fJz737nRg12Cz/wpZGSvOOO5isKI70i7vqfnfX175ILw4+81SCWrGnY3RKnaW2d7Zmuk3xPc4QRdPrZNIaYsBWilEbDmCWU1Wsl+YTJkkJGsEzqzkJBzPy6fQWOQAALgMjpAP54EvSobf81HSFC9HWF7ajgEvCM4ofjk+T+OP54Cun4MPperoAFBVvlI05J1rGXmWea6i1Ghjcbl31hKX4rfAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAB7bpwPoL02WGCCVhCsbDkq9GeFUwF01opVyFTijZlTUoTf5RcaR2qAH9/irkLjZeFeyozzC5mGvIVruBwnx/6l4PcAMxKK4YiheFVoO/dytpGMCj6ToNmKpjlXzOLAHelieWIUDtAFSYzENjIO01PyXTGYpxebpQCocJBvppj5HqARS9iNPcqBltMhxWrWmMu81tOG3Y7yd2xsIYXk6KjaoefLeN8Was4BPJ0zR6tTSEm6ZOvSRvlppqh84kz7LmWem7gGHAsY2G3tWBUmOdO/SMNMThqV62yLf7sKsuoE1w06lfmrf6D2zGwoEyz+TT6fdSkc34Yeh7+c01X6nFWU="],"x5t":"geiCPGtT_10T8xGLUK1LA0_YQEE","x5t#S256":"dLp3_QNGwMbYll5VecnR8Q9NSeFVfqJPBTa2_8qf48I"},{"kid":"tW6ae7TomE6_2jooM-sf9N_6lWg7HNtaQXrDsElBzM4","kty":"RSA","alg":"PS512","use":"sig","n":"p32N7jqKfMUB6_dKY1uZ3wizzPlBAXg9XrntfUcwNLRPfTBnshpt4uQBf3T8fexkbzhtR18oHvim-YvcWfC5eLGQmWHYiVwACa_C7oGqx51ijK2LRbUg4TKhnZX2X3Ld9xvr3HsosKh2UXn_Ay8nuvdfH-U6S7btT6a-AIFlt3BpqZP0EOl7rY-ie8nXoA13xX6BoyzYiNcugdYCU6czQcmTIJ1JLS0zohi4aTNehRt-1VMRpIMx7q7Ouq3Zhbi7RcDo-_D8FPRhWc2eEKd-h8ebFTIxEOrkguBIomjEFTf3SfYbOB_h-14v9Q2yz-NzyId3-ujRCQGC0hn-cixe2w","e":"AQAB","x5c":["MIICmzCCAYMCBgF4BKAxqzANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjEwMzA1MjMwMDEwWhcNMzEwMzA1MjMwMTUwWjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCnfY3uOop8xQHr90pjW5nfCLPM+UEBeD1eue19RzA0tE99MGeyGm3i5AF/dPx97GRvOG1HXyge+Kb5i9xZ8Ll4sZCZYdiJXAAJr8LugarHnWKMrYtFtSDhMqGdlfZfct33G+vceyiwqHZRef8DLye6918f5TpLtu1Ppr4AgWW3cGmpk/QQ6Xutj6J7ydegDXfFfoGjLNiI1y6B1gJTpzNByZMgnUktLTOiGLhpM16FG37VUxGkgzHurs66rdmFuLtFwOj78PwU9GFZzZ4Qp36Hx5sVMjEQ6uSC4EiiaMQVN/dJ9hs4H+H7Xi/1DbLP43PIh3f66NEJAYLSGf5yLF7bAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAHVWNBTExqlg4LTcyhUXI5U0iNPcMIVdKDoGPDc3EPjXyYNyjURX0oZ6b1Wv5t+XGmpZRqJNYb92xraQatIzLEsRn4IrmzViP+dIyFU8BEDubixTxeqx7LSw2j6LIFnZ05XdmWknlksNTlqi4CT6KL+1c24+QU3CcmU3mkQEIPA2yC4SdAB1oXI0jh49uP6a+JrE7JREZGAdwbIpZ1cqV6acPiJW3tOYfLrHwo7KYn3KwJvIBHXgFBNwx7fl2gYNQ0VEGKub3qVwW5RO5R/6Tcla9uZEfEiamms/Pn4hFA1qbsNHtA9IRGVRSmVeBKDxRvo0fxOUXp+NuZxEnhsoP3I="],"x5t":"f1l1fxICz1fe9mI-sSrtc19EDhU","x5t#S256":"NUJWRA4ADpLEg_SMkSoE4FKQN0H1Tlz85L-i7puVcqQ"},{"kid":"Lx1FmayP2YBtxaqS1SKJRJGiXRKnw2ov5WmYIMG-BLE","kty":"RSA","alg":"PS384","use":"sig","n":"q7WM4SnrdzlFSo_A1DRhc-8Ho-pBsfs49kGRbw3O_OKFIUyZrzHaRuovW_QaEAyiO3HX8CNcGPcpHdmpl4DhTGEBLcd6xXtCaa65ct00Mq7ZHCRRCrKLh6lJ0rY9fP8vCV0RBigpkNoRfrqLQQN4VeVFTbGSrDaS0LzPbap0-q5FKXUR-OQmQEtOupXhKFQtbB73tL83YnG6Swl7nXsx54ulEoDzcCCYt7pjCVVp7L9fzI2_ucTdtQclAJVQZGKpsx7vabOJuiMUwuAIz56lOJyXRMePsW8UogwC4FA2A52STsYlhOPsDEW4iIExFVNqs-CGoDGhYLIavaCkZhXM0w","e":"AQAB","x5c":["MIICmzCCAYMCBgF4HR+9XjANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjEwMzEwMTcxMDIyWhcNMzEwMzEwMTcxMjAyWjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCrtYzhKet3OUVKj8DUNGFz7wej6kGx+zj2QZFvDc784oUhTJmvMdpG6i9b9BoQDKI7cdfwI1wY9ykd2amXgOFMYQEtx3rFe0Jprrly3TQyrtkcJFEKsouHqUnStj18/y8JXREGKCmQ2hF+uotBA3hV5UVNsZKsNpLQvM9tqnT6rkUpdRH45CZAS066leEoVC1sHve0vzdicbpLCXudezHni6USgPNwIJi3umMJVWnsv1/Mjb+5xN21ByUAlVBkYqmzHu9ps4m6IxTC4AjPnqU4nJdEx4+xbxSiDALgUDYDnZJOxiWE4+wMRbiIgTEVU2qz4IagMaFgshq9oKRmFczTAgMBAAEwDQYJKoZIhvcNAQELBQADggEBADTgP3SrcG3p9XUB7sM4a2IeY0J4bSEtqlZBuHgdgekYJ5DXETJ3hV/82GjitU50NBup0IJyI9KZ0KCwqHIKC2Jn/6biOpM9Ipk4BtNVzx3qKNsDac9qZmyMpm4V9QuWakajknerhwyynG3siGUntbPmLvf5UKvKtbiKlWS4dBPwfedIUnC85mYEnNKSzSI1NiM6TWHB9zQYkARXlb89sh0HBYs08BfRMyBVM+l3OczIyGeQAfhcL+pxPP/0jqPr1ctHUBj2zXkjZxDw1oJFgeD9GDtPcjc3spB20vsRtQUBlzbJElbGflqWGHJK5l5n7gNd3ZXZT0HJ+wUpPE8EUaM="],"x5t":"fjRYR1986VCLzbaZaw5r25UKahw","x5t#S256":"ZHNHpizlsjD3qSZh7gJQQBu8W9jBL2HR0y7-3u2Wb-g"},{"kid":"gnmAfvmlsi3kKH3VlM1AJ85P2hekQ8ON_XvJqs3xPD8","kty":"RSA","alg":"RS384","use":"sig","n":"qUNQewKl3APQcbpACMNJ2XphPpupt395z6OZvj5CW9tiRXY3J7dqi8U0bWoIhtmmc7Js6hjp-A5W_FVStuXlT1hLyjJsHeu9ZVPnfIl2MnYN83zQBKw8E4mFsVv0UXNvkVPBF_k0yXrz-ABleWLOgFGnkNU9csc3Z5aihHcwRmC_oS7PZ9Vc-l0xBCyF3YRHI-al8ppSHwFreOweF3-JP3poNAXd906_tjX2KlHSJmNqcUNiSfEluyCp02ALlRFKXUQ1HlfSupHcHySDlanfUyIzZgM9ysCvC1vfNdAuwZ44oUBMul_XPxxhzlewL2Y8PtSDLUDWGTIou8M8049D8Q","e":"AQAB","x5c":["MIICmzCCAYMCBgF4BJVfaDANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjEwMzA1MjI0ODIxWhcNMzEwMzA1MjI1MDAxWjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCpQ1B7AqXcA9BxukAIw0nZemE+m6m3f3nPo5m+PkJb22JFdjcnt2qLxTRtagiG2aZzsmzqGOn4Dlb8VVK25eVPWEvKMmwd671lU+d8iXYydg3zfNAErDwTiYWxW/RRc2+RU8EX+TTJevP4AGV5Ys6AUaeQ1T1yxzdnlqKEdzBGYL+hLs9n1Vz6XTEELIXdhEcj5qXymlIfAWt47B4Xf4k/emg0Bd33Tr+2NfYqUdImY2pxQ2JJ8SW7IKnTYAuVEUpdRDUeV9K6kdwfJIOVqd9TIjNmAz3KwK8LW9810C7BnjihQEy6X9c/HGHOV7AvZjw+1IMtQNYZMii7wzzTj0PxAgMBAAEwDQYJKoZIhvcNAQELBQADggEBABoThxhMd7Xiq4x0GJeoJFv2yDKXCL3dJEAEWtOr2+PqdeJl/ZfOxBXynIvrdtYnQdICztN5ydEgDsZ02piDsxZ+s/0SA0iqjw/MEoBYobmr8V+xwUv+WtRLpTBXqWGMuG7NEtrbjKid0iKLLAOAU4dcHQ49iOF9VLnbTkf1EXp4iphJreaubOXMwT6/JDzQPT1dRR34hlhYeKKzMSA0Cz5aYL1tI+eH12rar0MDczXykLChNS/8MlyTzreEf0siUiS9S1kj/lOZKQDg9E/z8fm5vmHEHzAVwf4ON5iO29tDsqLw7BeJqC4AESjliXIqMrdpFynfPnIsGgf3dnph5BM="],"x5t":"CmRnQVduZWtEsdOC4mauUUsSWxA","x5t#S256":"BvC0LmuM8ZIApN3TQQZWWbGO-d082Ah5d3D6vPvahGw"},{"kid":"CGt0ZWS4Lc5faiKSdi0tU0fjCAdvGROQRGU9iR7tV0A","kty":"EC","alg":"ES256","use":"sig","crv":"P-256","x":"DPW7n9yjfE6Rt-VvVmEdeu4QdW44qifocAPPDxACDDY","y":"-ejsVw8222-hg2dJWx3QV0hE4-I0Ujp7ZsWebE68JE0"},{"kid":"C65q0EKQyhpd1m4fr7SKO2He_nAxgCtAdws64d2BLt8","kty":"RSA","alg":"RS256","use":"sig","n":"ja99ybDrLvw11Z4CvNlDI-kkqJEBpSnvDf0pZF2DvBlvYmeVYL_ChqIe8E9GyHUmLMdtO_jifSgOqE5b8vILwi1kZnJR7N857uEnbWM9YTeevi_RZ-E_hr4frW2NKJ78YGvCzwLKG2GgtSjj0zuTLnSaK8fCGzqXgy6paXNhgHUSZgGwvO0YItpMlyJeqEj1wGTWz1IyA1sguF1cC7K0fojPbPoBwrhvaAeoGRPLraE0rrBsQv8iiLwnRBIez9B1j0NiUG8Iad953Y7UzaKOAw8crIEK45NIK_yxHUpxqcHLjPIcRyIyJGioRyGK7cp-_7iPLOCutQc-u46mom1_ZQ","e":"AQAB","x5c":["MIICmzCCAYMCBgF4BJRpbzANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjEwMzA1MjI0NzE4WhcNMzEwMzA1MjI0ODU4WjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCNr33JsOsu/DXVngK82UMj6SSokQGlKe8N/SlkXYO8GW9iZ5Vgv8KGoh7wT0bIdSYsx207+OJ9KA6oTlvy8gvCLWRmclHs3znu4SdtYz1hN56+L9Fn4T+Gvh+tbY0onvxga8LPAsobYaC1KOPTO5MudJorx8IbOpeDLqlpc2GAdRJmAbC87Rgi2kyXIl6oSPXAZNbPUjIDWyC4XVwLsrR+iM9s+gHCuG9oB6gZE8utoTSusGxC/yKIvCdEEh7P0HWPQ2JQbwhp33ndjtTNoo4DDxysgQrjk0gr/LEdSnGpwcuM8hxHIjIkaKhHIYrtyn7/uI8s4K61Bz67jqaibX9lAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAHrGJFhVNiQupIwkn2jiW/jBobm9CHUxOwQL5E7WdRz5uaOJ0v62PrynOQE9xim9Qk8bT3q7DThZs66U9bpIk3msKVRgXRfn5FZy1H5RKOlEEFZhGakPqSlC1yPbhUNhHXMs3GTzdGMLtYaGvSy6XM/8/zqVqVwgh6BpbAR9RfiSdyaiNTSBriu+n/tHW934G9J8UIzdfpVcb0Yt9y4o0UgIXt64NtGFq7zmNJijH88AxBZFB6eUUmQQCczebzoAjyYbVOes5gGFzboVWcyLe3iyD0vvsAVHJViXeiGoxhpKnc8ryISpRUBzsKngf5uZo3bnrD9PHLYBoGOHgzII1xw="],"x5t":"5GNr3LeRXHWI4YR8-QTSsF98oTI","x5t#S256":"Dgd0_wZZqvRuf4GEISPNHREX-1ixTMIsrPeGzk0bCxs"}]}`)

// Create the JWKs from the resource at the given URL.
jwks, err := keyfunc.New(jwksJSON)
if err != nil {
	log.Fatalf("Failed to create JWKs from JSON.\nError: %s", err.Error())
}
```

Additional options can be passed to the [`keyfunc.Get`](https://pkg.go.dev/github.com/MicahParks/keyfunc#Get) function
via variadic arguments. See [`keyfunc.Options`](https://pkg.go.dev/github.com/MicahParks/keyfunc#Options) and the
additional features mentioned at the bottom of this `README.md`.

### Step 2: Use the [`keyfunc.JWKs`](https://pkg.go.dev/github.com/MicahParks/keyfunc#JWKs) 's [`JWKs.KeyFunc`](https://pkg.go.dev/github.com/MicahParks/keyfunc#JWKs.KeyFunc) method as the [`jwt.Keyfunc`](https://pkg.go.dev/github.com/dgrijalva/jwt-go@v3.2.0+incompatible#Keyfunc) when parsing tokens

```go
// Parse the JWT.
token, err := jwt.Parse(jwtB64, jwks.KeyFunc)
if err != nil {
	return nil, fmt.Errorf("failed to parse token: %w", err)
}
```

The [`JWKs.KeyFunc`](https://pkg.go.dev/github.com/MicahParks/keyfunc#JWKs.KeyFunc) method will automatically select the
key with the matching `kid` (if present) and return its public key as the correct Go type to its caller.

## Fork support

Some packages use forks of [github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go). This package aims to
support the most popular use cases of these forks.

If additional forks are required for your use case, please feel free to open an issue or PR.

### Support for [github.com/auth0/go-jwt-middleware](https://github.com/auth0/go-jwt-middleware)
The [github.com/auth0/go-jwt-middleware](https://github.com/auth0/go-jwt-middleware) package requires
the [github.com/form3tech-oss/jwt-go](https://github.com/form3tech-oss/jwt-go) fork. For this use case, use
the [`keyfunc.JWKs`](https://pkg.go.dev/github.com/MicahParks/keyfunc#JWKs) 's
[`JWKs.KeyFuncF3T`](https://pkg.go.dev/github.com/MicahParks/keyfunc#JWKs.KeyFuncF3T) method.

#### Example snippet
Please also see the `examples` directory.
```go
// Create the middleware provider.
jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{

	// Use the correct version of the KeyFunc method.
	ValidationKeyGetter: jwks.KeyFuncF3T,

	// Always ensure that you set your signing method to avoid tokens choosing the "none" method.
	//
	// This shouldn't matter for this keyfunc package, as the JWKs should be trusted and determines the key type,
	// but it's good practice.
	// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
	SigningMethod: jwt.SigningMethodRS256,
})
```

## Test coverage

Test coverage is currently at `85.3%`.

This is with current and expired JWTs, but the hard coded ones are now expired. Using non-expired JWTs would require
signing JWTs during testing and would allow for additional error checking. But a bit overkill since I've already done
that error checking when the JWTs were valid with no changes. A PR for this that does not introduce any dependencies is
welcome though.

## Additional features

* A background refresh of the JWKs keys can be performed. This is possible by passing
  [`keyfunc.Options`](https://pkg.go.dev/github.com/MicahParks/keyfunc#Options) via a variadic argument to the
  [`keyfunc.Get`](https://pkg.go.dev/github.com/MicahParks/keyfunc#Get) function.
	* A custom background refresh interval can be specified.
	* A custom background refresh request context timeout can be specified. Defaults to one minute.
	* A custom background refresh error handling function can be specified. If none is specified, errors go unhandled
	  silently.
* JWTs with a previously unseen `kid` can prompt an automatic refresh of the remote JWKs resource.
* A custom HTTP client can be used. This is possible by passing
  [`keyfunc.Options`](https://pkg.go.dev/github.com/MicahParks/keyfunc#Options) via a variadic argument to the
  [`keyfunc.Get`](https://pkg.go.dev/github.com/MicahParks/keyfunc#Get) function.

## TODO

- [ ] Add HMAC support?
