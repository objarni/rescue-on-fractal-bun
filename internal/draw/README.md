Features left TODO

- TextOp (worry: decide on anchor point placement)
- ScaleOp (Is this really needed? Menu screen can use larger font instead...)
- MirrorOp (at least along vertical axis)

/*
annoyances with draw.*
- color/colored
[x] color takes RGBA instead of Color interface
- Coordinate
- draw.
- intellisense confusion between structs and funcs
  (how avoid without hiding something from tests?)
  */

Potential improvements

This little 'internal library' called draw is starting to look good, however there are some warts I'd like to improve upon:

0. draw.Image could take only the "Image enum" instead
   of an map enum->*sprite, instead the map could be
   part of the Render(...) call!
1. The Coordinate vector never really 'flew'. Since the library is tightly integrated with Pixel anyway, just using pixel.Vec seems cleaner.
2. ImdOp Sequence uses array instead of combination through "Then" internally; this means annoyance w.r.t duplicated headBody code, used in other combinations. It's also incorrect, as it only inserts "  " in front of the first line and not others
3. Is it necessary to have Sequence twice...? It's on both
   ImdOp and WinOp.
4. WinOp Color resets to White; this is a bug since it
   means:
   
      ```
      Color rgb1 (
         (Color rgb2 (Image im1),
         Image im2)
      )
      ```
   
   will render im2 in white instead of rgb1.
   Not sure what is the solution though: sending not
   only Matrix but also Color recursively through
   Render is a solution, but feels over-the-top.

