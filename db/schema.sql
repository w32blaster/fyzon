PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS terms;
DROP TABLE IF EXISTS translations;

-- Project
CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

-- Language with codes and translations
CREATE TABLE terms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL,
    comment TEXT NOT NULL,
    project_id INTEGER,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE TABLE translations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    translation TEXT NOT NULL,
    language_code TEXT NOT NULL,
    is_default BOOLEAN DEFAULT false,
    term_id INTEGER,
    FOREIGN KEY (term_id) REFERENCES terms(id)
);
CREATE INDEX language_code_index ON translations(language_code);


-- Insert few Projects
INSERT INTO projects(id, name) values(1, "Project One");
INSERT INTO projects(id, name) values(2, "Project Two");


-- Insert some terms
INSERT INTO terms(id, code, comment, project_id) values(1, "main.page.title", "Top title on the main page", 1);
INSERT INTO terms(id, code, comment, project_id) values(2, "main.page.description", "Main page description", 1);

INSERT INTO terms(id, code, comment, project_id) values(3, "contact.form.title", "Contact form title", 2);
INSERT INTO terms(id, code, comment, project_id) values(4, "contact.form.submit", "The button on the form submit", 2);
INSERT INTO terms(id, code, comment, project_id) values(5, "contact.us.header", "Contact us page header", 2);

-- insert some translations
INSERT INTO translations(translation, language_code, term_id) values("The Title", "en", 1);
INSERT INTO translations(translation, language_code, term_id) values("Заголовок", "ru", 1);

INSERT INTO translations(translation, language_code, term_id) values("The description", "en", 2);
INSERT INTO translations(translation, language_code, term_id) values("Описание", "ru", 2);

INSERT INTO translations(translation, language_code, term_id) values("Fill The form to contact us", "en", 3);
INSERT INTO translations(translation, language_code, term_id) values("Заполните форму, чтобы связаться с нами", "ru", 3);

INSERT INTO translations(translation, language_code, term_id) values("Submit", "en", 4);
INSERT INTO translations(translation, language_code, term_id) values("Отправить", "ru", 4);

INSERT INTO translations(translation, language_code, term_id) values("Contact us", "en", 5);
INSERT INTO translations(translation, language_code, term_id) values("Свяжитесь с нами", "ru", 5);
