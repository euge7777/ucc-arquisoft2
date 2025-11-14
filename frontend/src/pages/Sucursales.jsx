import { useState, useEffect } from 'react';
import { MapPin, Phone, Mail, Clock, Check, Star, Target, Car, Accessibility, X } from 'lucide-react';
import { mockSucursales } from '../data/mockData';
import '../styles/Sucursales.css';

const Sucursales = () => {
    const [sucursales, setSucursales] = useState([]);
    const [loading, setLoading] = useState(true);
    const [sucursalSeleccionada, setSucursalSeleccionada] = useState(null);

    useEffect(() => {
        // Simular carga de API
        setTimeout(() => {
            setSucursales(mockSucursales);
            setLoading(false);
        }, 500);
    }, []);

    const handleVerMapa = (direccion) => {
        const query = encodeURIComponent(direccion);
        window.open(`https://www.google.com/maps/search/?api=1&query=${query}`, '_blank');
    };

    const handleLlamar = (telefono) => {
        window.location.href = `tel:${telefono}`;
    };

    const handleEmail = (email) => {
        window.location.href = `mailto:${email}`;
    };

    if (loading) {
        return (
            <div className="sucursales-container">
                <div className="loading-message">Cargando sucursales...</div>
            </div>
        );
    }

    return (
        <div className="sucursales-container">
            <div className="sucursales-header">
                <h1>Nuestras Sucursales</h1>
                <p>Encontrá la sucursal más cercana a vos</p>
            </div>

            <div className="sucursales-grid">
                {sucursales.map((sucursal) => (
                    <div
                        key={sucursal.id}
                        className={`sucursal-card ${sucursal.destacada ? 'destacada' : ''}`}
                    >
                        {sucursal.destacada && (
                            <div className="destacada-badge"><Star size={18} className="inline mr-2" /> Destacada</div>
                        )}

                        <div className="sucursal-imagen">
                            <img
                                src={sucursal.imagen}
                                alt={sucursal.nombre}
                                onError={(e) => {
                                    e.target.src = "https://via.placeholder.com/800x400?text=Gimnasio";
                                }}
                            />
                        </div>

                        <div className="sucursal-content">
                            <h2>{sucursal.nombre}</h2>

                            <div className="sucursal-info">
                                <div className="info-item">
                                    <span className="info-icon"><MapPin size={20} /></span>
                                    <span className="info-texto">{sucursal.direccion}</span>
                                </div>

                                <div className="info-item">
                                    <span className="info-icon"><Phone size={20} /></span>
                                    <span className="info-texto">{sucursal.telefono}</span>
                                </div>

                                <div className="info-item">
                                    <span className="info-icon"><Mail size={20} /></span>
                                    <span className="info-texto">{sucursal.email}</span>
                                </div>

                                <div className="info-item horarios">
                                    <span className="info-icon"><Clock size={20} /></span>
                                    <span className="info-texto">{sucursal.horarios}</span>
                                </div>
                            </div>

                            <div className="sucursal-servicios">
                                <h3>Servicios disponibles:</h3>
                                <div className="servicios-lista">
                                    {sucursal.servicios.map((servicio, index) => (
                                        <span key={index} className="servicio-tag">
                                            <Check size={16} className="inline mr-1" /> {servicio}
                                        </span>
                                    ))}
                                </div>
                            </div>

                            <div className="sucursal-acciones">
                                <button
                                    className="btn-ver-mapa"
                                    onClick={() => handleVerMapa(sucursal.direccion)}
                                >
                                    <MapPin size={18} className="inline mr-2" /> Ver en Mapa
                                </button>
                                <button
                                    className="btn-llamar"
                                    onClick={() => handleLlamar(sucursal.telefono)}
                                >
                                    <Phone size={18} className="inline mr-2" /> Llamar
                                </button>
                                <button
                                    className="btn-email"
                                    onClick={() => handleEmail(sucursal.email)}
                                >
                                    <Mail size={18} className="inline mr-2" /> Email
                                </button>
                            </div>

                            <button
                                className="btn-ver-detalle"
                                onClick={() => setSucursalSeleccionada(sucursal)}
                            >
                                Ver más información
                            </button>
                        </div>
                    </div>
                ))}
            </div>

            <div className="sucursales-info-adicional">
                <div className="info-box">
                    <h3><Target size={24} className="inline mr-2" /> ¿Primera vez?</h3>
                    <p>Visitá cualquiera de nuestras sucursales y obtené una clase de prueba gratuita</p>
                </div>
                <div className="info-box">
                    <h3><Car size={24} className="inline mr-2" /> Estacionamiento</h3>
                    <p>Todas nuestras sucursales cuentan con estacionamiento disponible</p>
                </div>
                <div className="info-box">
                    <h3><Accessibility size={24} className="inline mr-2" /> Accesibilidad</h3>
                    <p>Instalaciones completamente accesibles para personas con movilidad reducida</p>
                </div>
            </div>

            {/* Modal de detalle de sucursal */}
            {sucursalSeleccionada && (
                <div className="modal-overlay" onClick={() => setSucursalSeleccionada(null)}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <button
                            className="modal-close"
                            onClick={() => setSucursalSeleccionada(null)}
                        >
                            <X size={24} />
                        </button>

                        <div className="modal-header">
                            <h2>{sucursalSeleccionada.nombre}</h2>
                        </div>

                        <div className="modal-body">
                            <img
                                src={sucursalSeleccionada.imagen}
                                alt={sucursalSeleccionada.nombre}
                                className="modal-imagen"
                            />

                            <div className="modal-info">
                                <h3>Información de Contacto</h3>
                                <p><strong>Dirección:</strong> {sucursalSeleccionada.direccion}</p>
                                <p><strong>Teléfono:</strong> {sucursalSeleccionada.telefono}</p>
                                <p><strong>Email:</strong> {sucursalSeleccionada.email}</p>
                                <p><strong>Horarios:</strong> {sucursalSeleccionada.horarios}</p>
                            </div>

                            <div className="modal-servicios">
                                <h3>Servicios e Instalaciones</h3>
                                <ul>
                                    {sucursalSeleccionada.servicios.map((servicio, index) => (
                                        <li key={index}><Check size={18} className="inline mr-2" /> {servicio}</li>
                                    ))}
                                </ul>
                            </div>

                            <div className="modal-acciones">
                                <button
                                    className="btn-ver-mapa"
                                    onClick={() => handleVerMapa(sucursalSeleccionada.direccion)}
                                >
                                    <MapPin size={18} className="inline mr-2" /> Cómo Llegar
                                </button>
                                <button
                                    className="btn-contactar"
                                    onClick={() => handleLlamar(sucursalSeleccionada.telefono)}
                                >
                                    <Phone size={18} className="inline mr-2" /> Contactar
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Sucursales;
