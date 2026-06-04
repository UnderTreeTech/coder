{{- $short := (shortname .Name "err" "res" "sqlStr" "db" "XOLog") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- end -}}

/********************Generate By XO: write method segment*********************/
// Add{{ .Name }} add the {{ .Name }} to the database
func (d *dao) Add{{ .Name }}(ctx context.Context, {{$short}} *model.{{.Name}}) (err error) {
    // sql insert query, primary key provided by autoincrement
	sqlStr, args, err := squirrel.
        		Insert(model.Get{{  .Name  }}TableName()).
        		Columns(model.Get{{  .Name  }}FieldStringSlice()...).
        		Values(model.Get{{  .Name  }}AddField({{$short}})...).
        		PlaceholderFormat(d.PlaceHolder()).
        		ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
    sqlStr, err = drivers.QuoteSQL(d.driver, sqlStr)
    if err != nil {
       return
    }

    // run query
    log.Debug(ctx, "Add{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        _, err = d.db.Exec(ctx, sqlStr, args...)
    } else {
        _, err = tx.Exec(sqlStr, args...)
    }
    return
}

// Batch{{ .Name }} batch add {{ .Name }} to the database
func (d *dao) BatchAdd{{ .Name }}s(ctx context.Context, list []*model.{{.Name}}) (err error) {
    if 0 == len(list) {
        return
    }

    if 50 < len(list) {
        return errors.New("batch operation supports up to 50")
    }

    // sql insert query, primary key provided by autoincrement
	sqlBuilder := squirrel.
        Insert(model.Get{{  .Name  }}TableName()).
        Columns(model.Get{{  .Name  }}FieldStringSlice()...)

    for _, {{ $short }} := range list {
        sqlBuilder = sqlBuilder.Values(model.Get{{  .Name  }}AddField({{$short}})...)
    }

    sqlStr, args, err := sqlBuilder.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return err
    }

    // parse sql to adapter databases
    sqlStr, err = drivers.QuoteSQL(d.driver, sqlStr)
    if err != nil {
       return
    }

    // run query
    log.Debug(ctx, "BatchAdd{{ .Name }}s", log.String("sql", fmt.Sprint(sqlStr, args)))
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        _, err = d.db.Exec(ctx, sqlStr, args...)
    } else {
        _, err = tx.Exec(sqlStr, args...)
    }
    return
}

// Edit{{ .Name }} edit {{ .Name }} in the database
func (d *dao) Edit{{ .Name }}(ctx context.Context, setMap map[string]interface{}, condition map[string]interface{}) (err error) {
    // build sql
	sqlStr, args, err := squirrel.
    		Update(model.Get{{  .Name  }}TableName()).
    		SetMap(setMap).
    		Where(condition).
    		PlaceholderFormat(d.PlaceHolder()).
    		ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
    sqlStr, err = drivers.QuoteSQL(d.driver, sqlStr)
    if err != nil {
       return
    }

	// run query
	log.Debug(ctx, "Edit{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
	if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        _, err = d.db.Exec(ctx, sqlStr, args...)
    } else {
        _, err = tx.Exec(sqlStr, args...)
    }
    return
}

// Delete{{ .Name }} delete {{ .Name }} from the database
func (d *dao)  Delete{{ .Name }}(ctx context.Context,condition map[string]interface{}) (err error) {
    sqlStr,args, err := squirrel.
           Delete(model.Get{{  .Name  }}TableName()).
           Where(condition).
           PlaceholderFormat(d.PlaceHolder()).
           ToSql()
    if err != nil {
            return
    }

    // parse sql to adapter databases
    sqlStr, err = drivers.QuoteSQL(d.driver, sqlStr)
    if err != nil {
       return
    }

    // run query
    log.Debug(ctx, "Delete{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        _, err = d.db.Exec(ctx, sqlStr, args...)
    } else {
        _, err = tx.Exec(sqlStr, args...)
    }
    return
}

/******************Generate By XO: read method segment**********************/
// Find{{ .Name }} query one {{ .Name }} by condition from the database
func (d *dao)  Find{{ .Name }}(ctx context.Context, condition map[string]interface{}) (res *model.{{ .Name }}, err error) {
    // init build
    build := squirrel.Select(model.Get{{ .Name }}FieldString()).From(model.Get{{  .Name  }}TableName())

    // forced master database query
    force := false
    if _, ok := condition["_force"]; ok {
        force = true
        delete(condition, "_force")
    }

    // build sql
    build, err = d.Analytic(build, condition)
    if err != nil {
        return
    }
    sqlStr, args, err := build.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
    sqlStr, err = drivers.QuoteSQL(d.driver, sqlStr)
    if err != nil {
       return
    }

    // run query
    log.Debug(ctx, "Find{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
    res = &model.{{ .Name }}{}
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        if force {
            err = d.db.Master().QueryRow(ctx, sqlStr, args...).Scan(model.Get{{ .Name }}ScanField(res)...)
        } else {
            err = d.db.QueryRow(ctx, sqlStr, args...).Scan(model.Get{{ .Name }}ScanField(res)...)
        }
    } else {
        err = tx.QueryRow(sqlStr, args...).Scan(model.Get{{ .Name }}ScanField(res)...)
    }
    return
}

// Find{{ .Name }}s all {{ .Name }} by condition from the database
func (d *dao)  Find{{ .Name }}s(ctx context.Context, condition map[string]interface{}) (res []*model.{{ .Name }}, err error) {
    // init build
    build := squirrel.Select(model.Get{{ .Name }}FieldString()).From(model.Get{{  .Name  }}TableName())

    // forced master database query
    force := false
    if _, ok := condition["_force"]; ok {
        force = true
        delete(condition, "_force")
    }

    // build sql
    build, err = d.Analytic(build, condition)
    if err != nil {
        return
    }
    sqlStr, args, err := build.PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
    sqlStr, err = drivers.QuoteSQL(d.driver, sqlStr)
    if err != nil {
       return
    }

    // run query
    log.Debug(ctx, "Find{{ .Name }}s", log.String("sql", fmt.Sprint(sqlStr, args)))
    var rows *sql.Rows
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        if force {
            rows, err = d.db.Master().Query(ctx, sqlStr, args...)
        } else {
            rows, err = d.db.Query(ctx, sqlStr, args...)
        }
    } else {
        rows, err = tx.Query(sqlStr, args...)
    }

    if err != nil {
        return
    }
    defer rows.Close()

    // load results
    res = make([]*model.{{ .Name }}, 0)
    for rows.Next() {
        // scan rows
        {{ $short }} := &model.{{ .Name }}{}
        err = rows.Scan(model.Get{{ .Name }}ScanField({{ $short }})...)
        if err != nil {
            return
        }

        res = append(res, {{ $short }})
    }
    return
}

// Count{{ .Name }} return count by condition from the database
func (d *dao)  Count{{ .Name }}(ctx context.Context, condition map[string]interface{}) (num int, err error) {
    // init build
    build := squirrel.Select("count(*)").From(model.Get{{  .Name  }}TableName())

    sqlStr, args, err := build.Where(condition).PlaceholderFormat(d.PlaceHolder()).ToSql()
    if err != nil {
        return
    }

    // parse sql to adapter databases
    sqlStr, err = drivers.QuoteSQL(d.driver, sqlStr)
    if err != nil {
       return
    }

    // run query
    log.Debug(ctx, "Count{{ .Name }}", log.String("sql", fmt.Sprint(sqlStr, args)))
    if tx, txErr := d.GetTxFromCtx(ctx); txErr != nil {
        err = d.db.QueryRow(ctx, sqlStr, args...).Scan(&num)
    } else {
        err = tx.QueryRow(sqlStr, args...).Scan(&num)
    }
    return
}
