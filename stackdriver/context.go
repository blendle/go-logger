// nolint:golint
package stackdriver

// MIT License
//
// Copyright (c) 2017 Tommy Chen
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

import (
	"go.uber.org/zap/zapcore"
)

// The schema is based on: https://cloud.google.com/error-reporting/docs/formatting-error-messages

type ServiceContext struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

func (s *ServiceContext) Clone() *ServiceContext {
	return &ServiceContext{
		Service: s.Service,
		Version: s.Version,
	}
}

func (s *ServiceContext) MarshalLogObject(e zapcore.ObjectEncoder) error {
	e.AddString("service", s.Service)
	e.AddString("version", s.Version)
	return nil
}

type Context struct {
	User           string          `json:"user"`
	HTTPRequest    *HTTPRequest    `json:"httpRequest"`
	ReportLocation *ReportLocation `json:"reportLocation"`
}

func (c *Context) Clone() *Context {
	output := &Context{
		User: c.User,
	}

	if c.HTTPRequest != nil {
		output.HTTPRequest = c.HTTPRequest.Clone()
	}

	if c.ReportLocation != nil {
		output.ReportLocation = c.ReportLocation.Clone()
	}

	return output
}

func (c *Context) MarshalLogObject(e zapcore.ObjectEncoder) error {
	var err error

	e.AddString("user", c.User)

	if c.HTTPRequest != nil {
		if err = e.AddObject("httpRequest", c.HTTPRequest); err != nil {
			return err
		}
	}

	if c.ReportLocation != nil {
		if err = e.AddObject("reportLocation", c.ReportLocation); err != nil {
			return err
		}
	}

	return err
}

type HTTPRequest struct {
	Method             string `json:"method"`
	URL                string `json:"url"`
	UserAgent          string `json:"userAgent"`
	Referrer           string `json:"referrer"`
	ResponseStatusCode int    `json:"responseStatusCode"`
	RemoteIP           string `json:"remoteIp"`
}

func (h *HTTPRequest) Clone() *HTTPRequest {
	return &HTTPRequest{
		Method:             h.Method,
		URL:                h.URL,
		UserAgent:          h.UserAgent,
		Referrer:           h.Referrer,
		ResponseStatusCode: h.ResponseStatusCode,
		RemoteIP:           h.RemoteIP,
	}
}

func (h *HTTPRequest) MarshalLogObject(e zapcore.ObjectEncoder) error {
	e.AddString("method", h.Method)
	e.AddString("url", h.URL)
	e.AddString("userAgent", h.UserAgent)
	e.AddString("referrer", h.Referrer)
	e.AddInt("responseStatusCode", h.ResponseStatusCode)
	e.AddString("remoteIp", h.RemoteIP)
	return nil
}

type ReportLocation struct {
	FilePath     string
	LineNumber   int
	FunctionName string
}

func (r *ReportLocation) Clone() *ReportLocation {
	return &ReportLocation{
		FilePath:     r.FilePath,
		LineNumber:   r.LineNumber,
		FunctionName: r.FunctionName,
	}
}

func (r *ReportLocation) MarshalLogObject(e zapcore.ObjectEncoder) error {
	e.AddString("filePath", r.FilePath)
	e.AddInt("lineNumber", r.LineNumber)
	e.AddString("functionName", r.FunctionName)
	return nil
}
