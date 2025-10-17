-- REGISTROS PARA TESTEAR LA BASE DE DATOS

-- Usuarios
INSERT INTO usuarios (id_usuario, nombre, apellido, username, password, is_admin) VALUES
  (2, 'Ana',    'Gómez',    'ana.gomez',    'pass123', 0),
  (3, 'Bruno',  'Pérez',    'bruno.pere',   'abc321',  0),
  (4, 'Carla',  'López',    'carla.lopez',  'xyz789',  1),
  (5, 'David',  'Santos',   'david.santos', 'qwe456',  0),
  (6, 'Elena',  'Martín',   'elena.mart',   'zxc852',  0)
;

UPDATE usuarios SET password = sha2(password, 256) WHERE id_usuario BETWEEN 2 AND 6;

-- Actividades
INSERT INTO actividads (id_actividad, foto_url, titulo, descripcion, cupo, dia, horario_inicio, horario_final, instructor, categoria) VALUES
  (1, 'https://images.unsplash.com/photo-1554302242-40743152783a?q=80&w=1886&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Karate Miyagi-do',        'Clase de karate con el estilo del señor Miyagi', 10, 'Sabado',    '2025-06-06 07:00:00', '2025-06-06 09:00:00', 'Sr. Miyagi',      'karate'),
  (2, 'https://images.unsplash.com/photo-1552196563-55cd4e45efb3?q=80&w=1926&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Yoga Suave',              'Clase de yoga para principiantes',               15, 'Lunes',     '2025-06-06 10:00:00', '2025-06-06 11:00:00', 'Laura Ruiz',      'yoga'),
  (3, 'https://images.unsplash.com/photo-1534258936925-c58bed479fcb?q=80&w=1931&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Crossfit Básico',         'Entrenamiento funcional de fuerza',              10, 'Miercoles', '2025-06-06 18:30:00', '2025-06-06 19:30:00', 'Martín Díaz',     'fitness'),
  (4, 'https://images.unsplash.com/photo-1530549387789-4c1017266635?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Natación Adultos',        'Lecciones de natación nivel intermedio',         8,  'Viernes',   '2025-06-06 09:00:00', '2025-06-06 10:00:00', 'Sandra Pérez',    'natación'),
  (5, 'https://images.unsplash.com/photo-1717500251716-27057c48ace4?q=80&w=1887&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Pilates',                 'Pilates para tonificar el core',                 12, 'Martes',    '2025-06-06 17:00:00', '2025-06-06 18:00:00', 'Carlos Méndez',   'pilates'),
  (6, 'https://images.unsplash.com/photo-1554470166-20d3f466089b?q=80&w=1886&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Spinning',                'Clase de ciclismo indoor de alta intensidad',    20, 'Jueves',    '2025-06-06 19:00:00', '2025-06-06 20:00:00', 'Lucía Herrera',   'ciclismo'),
  (7, 'https://images.unsplash.com/photo-1622599511051-16f55a1234d0?q=80&w=1936&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Boxeo profesional',       'Clase de boxeo',                                 4,  'Martes',    '2025-06-06 19:00:00', '2025-06-06 20:00:00', 'Mike Tyson',      'boxeo'),
  (8, 'https://images.unsplash.com/photo-1534438327276-14e5300c3a48?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Entrenamiento de fuerza', 'Entrenamiento con pesas y maquinas',             12, 'Miercoles', '2025-06-06 19:00:00', '2025-06-06 20:00:00', 'Anatoli Cleaner', 'fuerza'),
  (9, 'https://images.unsplash.com/photo-1616279969856-759f316a5ac1?q=80&w=1965&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D', 'Entrenamiento de fuerza', 'Entrenamiento con pesas y maquinas',             12, 'Miercoles', '2025-06-06 20:00:00', '2025-06-06 21:00:00', 'Anatoli Cleaner', 'fuerza')
;

-- Inscripciones
INSERT INTO inscripcions (id_usuario, id_actividad) VALUES
  (2, 1),
  (2, 7),
  (3, 1),
  (3, 7),
  (6, 1),
  (6, 2),
  (6, 4),
  (6, 5),
  (6, 7)
;

-- PARA ELIMINAR TODAS LAS TABLAS
-- set FOREIGN_KEY_CHECKS = 0;
-- drop table inscripcions;
-- drop table actividads;
-- drop table usuarios;
-- drop view actividads_lugares;
-- set FOREIGN_KEY_CHECKS = 1;

-- PARA BORRAR TODOS LOS REGISTROS
-- set FOREIGN_KEY_CHECKS = 0;
-- delete from usuarios;
-- delete from actividads;
-- delete from inscripcions;
-- set FOREIGN_KEY_CHECKS = 1;
