package repository

import "context"

func (r *Repository) FindByPhone(ctx context.Context, input FindByPhoneInput) (FindByPhoneOutput, error) {
	stmt, err := r.Db.PrepareContext(ctx, `SELECT id, slug, full_name, phone, password FROM users where phone=?`)
	if nil != err {
		return FindByPhoneOutput{}, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	var output FindByPhoneOutput
	if err = stmt.QueryRow(
		input.Phone,
	).Scan(
		&output.Id,
		&output.Slug,
		&output.FullName,
		&output.Phone,
		&output.Password,
	); nil != err {
		return FindByPhoneOutput{}, err
	}

	return output, nil
}

func (r *Repository) FindBySlug(ctx context.Context, input FindBySlugInput) (FindBySlugOutput, error) {
	stmt, err := r.Db.PrepareContext(ctx, `SELECT slug, full_name, phone, password FROM users where slug=?`)
	if nil != err {
		return FindBySlugOutput{}, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	var output FindBySlugOutput
	if err = stmt.QueryRow(
		input.Slug,
	).Scan(
		&output.Slug,
		&output.FullName,
		&output.Phone,
		&output.Password,
	); nil != err {
		return FindBySlugOutput{}, err
	}

	return output, nil
}

func (r *Repository) Store(ctx context.Context, input RegistrationInput) (RegistrationOutput, error) {
	stmt, err := r.Db.PrepareContext(ctx, `INSERT INTO users (slug, full_name, phone, password) VALUES (?, ?, ?, ?) RETURNING id`)
	if nil != err {
		return RegistrationOutput{}, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	var output RegistrationOutput
	if err = stmt.QueryRow(
		input.Slug,
		input.FullName,
		input.Phone,
		input.Password,
	).Scan(&output.Id); nil != err {
		return RegistrationOutput{}, err
	}

	return output, nil
}

func (r *Repository) Put(ctx context.Context, input UpdateUserInput) error {
	stmt, err := r.Db.PrepareContext(ctx, `UPDATE users SET full_name=?, phone=? where slug=?`)
	if nil != err {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(input.FullName, input.Phone, input.Slug)
	if nil != err {
		return err
	}

	return nil
}
