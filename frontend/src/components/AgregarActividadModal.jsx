import React, { useState } from 'react';
import '../styles/Modal.css';

const AgregarActividadModal = ({ onClose, onSave }) => {
    const [formData, setFormData] = useState({
        titulo: '',
        descripcion: '',
        cupo: '',
        dia: '',
        hora_inicio: '',
        hora_fin: '',
        foto_url: '',
        instructor: '',
        categoria: ''
    });
    const [error, setError] = useState('');
    const [validationErrors, setValidationErrors] = useState({});

    const validateForm = () => {
        const errors = {};
        
        if (!formData.titulo.trim()) {
            errors.titulo = 'El título es requerido';
        } else if (formData.titulo.length < 3) {
            errors.titulo = 'El título debe tener al menos 3 caracteres';
        }

        if (!formData.descripcion.trim()) {
            errors.descripcion = 'La descripción es requerida';
        }

        if (!formData.cupo || parseInt(formData.cupo) < 1) {
            errors.cupo = 'El cupo debe ser mayor a 0';
        }

        if (!formData.dia) {
            errors.dia = 'El día es requerido';
        }

        if (!formData.hora_inicio) {
            errors.hora_inicio = 'La hora de inicio es requerida';
        }

        if (!formData.hora_fin) {
            errors.hora_fin = 'La hora de fin es requerida';
        } else if (formData.hora_fin <= formData.hora_inicio) {
            errors.hora_fin = 'La hora de fin debe ser posterior a la hora de inicio';
        }

        if (!formData.instructor.trim()) {
            errors.instructor = 'El instructor es requerido';
        }

        if (!formData.categoria.trim()) {
            errors.categoria = 'La categoría es requerida';
        }

        setValidationErrors(errors);
        return Object.keys(errors).length === 0;
    };

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
        // Limpiar error de validación cuando el usuario modifica el campo
        if (validationErrors[name]) {
            setValidationErrors(prev => ({
                ...prev,
                [name]: ''
            }));
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!validateForm()) {
            return;
        }

        try {
            const token = localStorage.getItem('access_token');
            if (!token) {
                setError('No hay sesión activa. Por favor, inicie sesión nuevamente.');
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
                return;
            }

            const dataToSend = {
                ...formData,
                cupo: parseInt(formData.cupo, 10),
                dia: formData.dia.normalize("NFD").replace(/[\u0300-\u036f]/g, ""), // Eliminar acentos
                hora_inicio: formData.hora_inicio,
                hora_fin: formData.hora_fin
            };

            const response = await fetch('http://localhost:8080/actividades', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(dataToSend)
            });

            if (response.status === 401) {
                setError('Su sesión ha expirado. Por favor, inicie sesión nuevamente.');
                localStorage.removeItem('access_token');
                setTimeout(() => {
                    window.location.href = '/login';
                }, 2000);
                return;
            }

            if (response.ok) {
                onSave();
                onClose();
            } else {
                const errorData = await response.json();
                setError(errorData.error || 'Error al crear la actividad');
            }
        } catch (error) {
            console.error('Error:', error);
            setError('Error al conectar con el servidor');
        }
    };

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h2>Agregar Nueva Actividad</h2>
                {error && <div className="error-message">{error}</div>}
                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="titulo">Título:</label>
                        <input
                            type="text"
                            id="titulo"
                            name="titulo"
                            value={formData.titulo}
                            onChange={handleChange}
                            required
                        />
                        {validationErrors.titulo && <span className="error-text">{validationErrors.titulo}</span>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="descripcion">Descripción:</label>
                        <textarea
                            id="descripcion"
                            name="descripcion"
                            value={formData.descripcion}
                            onChange={handleChange}
                            required
                        />
                        {validationErrors.descripcion && <span className="error-text">{validationErrors.descripcion}</span>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="cupo">Cupo:</label>
                        <input
                            type="number"
                            id="cupo"
                            name="cupo"
                            value={formData.cupo}
                            onChange={handleChange}
                            required
                            min="1"
                        />
                        {validationErrors.cupo && <span className="error-text">{validationErrors.cupo}</span>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="dia">Día:</label>
                        <select
                            id="dia"
                            name="dia"
                            value={formData.dia}
                            onChange={handleChange}
                            required
                        >
                            <option value="">Seleccione un día</option>
                            <option value="Lunes">Lunes</option>
                            <option value="Martes">Martes</option>
                            <option value="Miercoles">Miercoles</option>
                            <option value="Jueves">Jueves</option>
                            <option value="Viernes">Viernes</option>
                            <option value="Sabado">Sabado</option>
                        </select>
                        {validationErrors.dia && <span className="error-text">{validationErrors.dia}</span>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="hora_inicio">Hora de inicio:</label>
                        <input
                            type="time"
                            id="hora_inicio"
                            name="hora_inicio"
                            value={formData.hora_inicio}
                            onChange={handleChange}
                            required
                        />
                        {validationErrors.hora_inicio && <span className="error-text">{validationErrors.hora_inicio}</span>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="hora_fin">Hora de fin:</label>
                        <input
                            type="time"
                            id="hora_fin"
                            name="hora_fin"
                            value={formData.hora_fin}
                            onChange={handleChange}
                            required
                        />
                        {validationErrors.hora_fin && <span className="error-text">{validationErrors.hora_fin}</span>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="foto_url">URL de la foto:</label>
                        <input
                            type="text"
                            id="foto_url"
                            name="foto_url"
                            value={formData.foto_url}
                            onChange={handleChange}
                        />
                        {validationErrors.foto_url && <span className="error-text">{validationErrors.foto_url}</span>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="instructor">Instructor:</label>
                        <input
                            type="text"
                            id="instructor"
                            name="instructor"
                            value={formData.instructor}
                            onChange={handleChange}
                            required
                        />
                        {validationErrors.instructor && <span className="error-text">{validationErrors.instructor}</span>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="categoria">Categoría:</label>
                        <input
                            type="text"
                            id="categoria"
                            name="categoria"
                            value={formData.categoria}
                            onChange={handleChange}
                            required
                        />
                        {validationErrors.categoria && <span className="error-text">{validationErrors.categoria}</span>}
                    </div>

                    <div className="form-buttons">
                        <button type="submit" className="btn-guardar">Crear Actividad</button>
                        <button type="button" className="btn-cancelar" onClick={onClose}>
                            Cancelar
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default AgregarActividadModal; 