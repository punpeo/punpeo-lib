package {{.PkgName}}

import (
    "net/http"
    {{if .HasRequest}}"gitlab.jianzhikeji.com/jz-backend/go-lib/zeropro/requests"{{end}}
    "gitlab.jianzhikeji.com/jz-backend/go-lib/zeropro/httpxpro/render"

   {{.ImportPackages}}
)

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        {{if .HasRequest}}var req types.{{.RequestType}}
        if err := requests.ParseAndValidate(r, &req); err != nil {
            render.JsonErrorResult(r, w, err)
            return
        }{{end}}



        l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
        {{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})

        {{if .HasResp}}render.JsonResult(r, w, resp, err){{else}}result.HttpResult(r, w, nil, err){{end}}
    }
}