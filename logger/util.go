package logger

func ParseMod(m string) (ModeOut, error) {
  v, ok := ModeKeyVal[m]
  if ok {
    return v, nil
  }

  return 0, ErrParseMode
}
