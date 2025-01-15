package glog

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type styleConfig struct {
  level map[log.Level]lipgloss.Style 
}

type Style func(*styleConfig) error;

func newStyle(styles ...Style) (*styleConfig, error) {
  stylecfg := styleConfig{
    level: log.DefaultStyles().Levels,
  }

  for _, styleOpt := range styles {
    if err := styleOpt(&stylecfg); err != nil {
      return nil, err
    }
  }

  return &stylecfg, nil
}
