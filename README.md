# distill-bin

Binary distillation shortcut and McCabe-Thiele stage rating. Given feed rate
and composition, distillate and bottoms targets, feed thermal quality `q` and
relative volatility (or a discrete equilibrium curve), it computes the minimum
reflux by Underwood, the minimum theoretical stages by Fenske, and — at a user
specified reflux — steps off the rectifying and stripping operating lines to
count theoretical stages including the reboiler.

Material balance is closed: `F = D + B` and `F zF = D xD + B xB` with the
residual reported. Equilibrium is the constant-volatility relation
`y = αx/(1+(α−1)x)`. The Gilliland correlation used is Eduljee's
`(N − Nmin)/(N + 1) = 0.75 (1 − x^0.5668)` with `x = (R − Rmin)/(R + 1)`.

## Build

```bash
go build .
go test ./...
```

## Usage

```bash
go run . rate example/benzene.json
go run . gilliland example/benzene.json
go run . stage example/benzene.json
cat example/benzene.json | go run . rate
```

Input JSON:

```json
{
  "feed": 100,
  "feed_composition": 0.45,
  "distillate_composition": 0.97,
  "bottoms_composition": 0.03,
  "alpha": 2.5,
  "q": 1.0,
  "reflux": 2.5
}
```

`rate` prints D, B, Rmin, Fenske Nmin, the stage count at the given reflux and
the tray-by-tray liquid composition. `gilliland` prints Nmin/Rmin and Gilliland
stage estimates at several reflux ratios. `stage` prints the full McCabe-Thiele
composition profile with feed tray.

## Validation

`xD`, `xB`, `zF` must lie in (0,1) with `xD > zF > xB`; `α <= 1` is reported as
"cannot separate"; a reflux below `Rmin` is reported as "insufficient reflux"
instead of looping forever. Missing feed rate, NaN, unknown errors exit
non-zero with a message on stderr.
