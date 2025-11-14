import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Waves, ClipboardList, Dumbbell, CreditCard, Calendar, AlertCircle, CheckCircle2, Clock, MapPin } from 'lucide-react';
import { getMockSuscripcionByUserId } from '../data/mockData';
import { ACTIVITIES_API, PAYMENTS_API } from '../config/api';
import '../styles/Dashboard.css';

const Dashboard = () => {
    const [suscripcion, setSuscripcion] = useState(null);
    const [inscripciones, setInscripciones] = useState([]);
    const [pagosRecientes, setPagosRecientes] = useState([]);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();
    const userId = localStorage.getItem("idUsuario");
    const username = localStorage.getItem("username") || "Usuario";

    useEffect(() => {
        fetchDashboardData();
    }, [userId]);

    const fetchDashboardData = async () => {
        try {
            setLoading(true);

            // Obtener suscripción (mock)
            const subData = getMockSuscripcionByUserId(userId);
            setSuscripcion(subData);

            // Obtener inscripciones (activities-api)
            try {
                const inscResponse = await fetch(ACTIVITIES_API.inscripciones, {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('access_token')}`
                    }
                });
                if (inscResponse.ok) {
                    const inscData = await inscResponse.json();
                    setInscripciones(inscData.filter(i => i.is_activa));
                }
            } catch (error) {
                console.error("Error al cargar inscripciones:", error);
            }

            // Obtener pagos recientes (payments-api)
            try {
                const pagosResponse = await fetch(PAYMENTS_API.paymentsByUser(userId), {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('access_token')}`
                    }
                });
                if (pagosResponse.ok) {
                    const pagosData = await pagosResponse.json();
                    // Tomar solo los últimos 3 pagos
                    setPagosRecientes(pagosData.slice(0, 3));
                }
            } catch (error) {
                console.error("Error al cargar pagos:", error);
            }

        } catch (error) {
            console.error("Error al cargar dashboard:", error);
        } finally {
            setLoading(false);
        }
    };

    const getDiasRestantes = (fechaVencimiento) => {
        const hoy = new Date();
        const vencimiento = new Date(fechaVencimiento);
        const diffTime = vencimiento - hoy;
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        return diffDays;
    };

    if (loading) {
        return (
            <div className="dashboard-container">
                <div className="loading-message">Cargando dashboard...</div>
            </div>
        );
    }

    const diasRestantes = suscripcion ? getDiasRestantes(suscripcion.fecha_vencimiento) : 0;
    const proximoVencer = diasRestantes <= 7 && diasRestantes > 0;

    return (
        <div className="dashboard-container">
            <div className="dashboard-header">
                <h1>
                    <Waves size={32} className="inline mr-2" style={{verticalAlign: 'middle'}} />
                    Bienvenido, {username}
                </h1>
                <p>Aquí está el resumen de tu actividad en GymPro</p>
            </div>

            <div className="dashboard-stats">
                <div className="stat-card">
                    <div className="stat-icon"><ClipboardList size={24} /></div>
                    <div className="stat-info">
                        <span className="stat-value">
                            {suscripcion ? (suscripcion.estado === 'activa' ? 'Activa' : 'Inactiva') : 'Sin plan'}
                        </span>
                        <span className="stat-label">Suscripción</span>
                    </div>
                </div>

                <div className="stat-card">
                    <div className="stat-icon"><Dumbbell size={24} /></div>
                    <div className="stat-info">
                        <span className="stat-value">{inscripciones.length}</span>
                        <span className="stat-label">Actividades Inscritas</span>
                    </div>
                </div>

                <div className="stat-card">
                    <div className="stat-icon"><CreditCard size={24} /></div>
                    <div className="stat-info">
                        <span className="stat-value">{pagosRecientes.length}</span>
                        <span className="stat-label">Pagos Recientes</span>
                    </div>
                </div>

                {suscripcion && (
                    <div className="stat-card">
                        <div className="stat-icon"><Calendar size={24} /></div>
                        <div className="stat-info">
                            <span className={`stat-value ${proximoVencer ? 'stat-warning' : ''}`}>
                                {diasRestantes > 0 ? `${diasRestantes} días` : 'Vencida'}
                            </span>
                            <span className="stat-label">Tiempo Restante</span>
                        </div>
                    </div>
                )}
            </div>

            <div className="dashboard-content">
                {/* Sección de Suscripción */}
                <div className="dashboard-section suscripcion-section">
                    <div className="section-header">
                        <h2>Mi Suscripción</h2>
                        <button
                            className="btn-ver-mas"
                            onClick={() => navigate('/mi-suscripcion')}
                        >
                            Ver detalle →
                        </button>
                    </div>

                    {suscripcion ? (
                        <div className="suscripcion-widget" style={{ borderLeftColor: suscripcion.plan?.color }}>
                            <div className="suscripcion-info">
                                <h3>{suscripcion.plan?.nombre}</h3>
                                <p>{suscripcion.plan?.descripcion}</p>
                                <div className="suscripcion-detalles-mini">
                                    <span>Vence: {new Date(suscripcion.fecha_vencimiento).toLocaleDateString('es-AR')}</span>
                                    <span className={`estado-mini ${suscripcion.estado === 'activa' ? 'activa' : ''}`}>
                                        {suscripcion.estado}
                                    </span>
                                </div>
                            </div>
                            {proximoVencer && (
                                <div className="alerta-mini">
                                    <AlertCircle size={18} className="inline mr-2" />
                                    Tu suscripción vence pronto
                                    <button
                                        className="btn-renovar-mini"
                                        onClick={() => navigate('/planes')}
                                    >
                                        Renovar
                                    </button>
                                </div>
                            )}
                        </div>
                    ) : (
                        <div className="no-suscripcion-widget">
                            <p>No tenés una suscripción activa</p>
                            <button
                                className="btn-suscribirse"
                                onClick={() => navigate('/planes')}
                            >
                                Ver Planes
                            </button>
                        </div>
                    )}
                </div>

                {/* Sección de Actividades */}
                <div className="dashboard-section actividades-section">
                    <div className="section-header">
                        <h2>Mis Actividades</h2>
                        <button
                            className="btn-ver-mas"
                            onClick={() => navigate('/actividades')}
                        >
                            Ver todas →
                        </button>
                    </div>

                    {inscripciones.length > 0 ? (
                        <div className="actividades-lista-mini">
                            {inscripciones.slice(0, 3).map((insc) => (
                                <div key={insc.id_inscripcion} className="actividad-mini">
                                    <div className="actividad-mini-info">
                                        <span className="actividad-icono"><Dumbbell size={20} /></span>
                                        <div>
                                            <span className="actividad-nombre">Actividad #{insc.id_actividad}</span>
                                            <span className="actividad-fecha">
                                                Inscrito: {new Date(insc.fecha_inscripcion).toLocaleDateString('es-AR')}
                                            </span>
                                        </div>
                                    </div>
                                    <span className="actividad-estado"><CheckCircle2 size={16} className="inline mr-1" /> Activa</span>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div className="no-actividades-widget">
                            <p>No estás inscrito a ninguna actividad</p>
                            <button
                                className="btn-inscribirse"
                                onClick={() => navigate('/actividades')}
                            >
                                Ver Actividades
                            </button>
                        </div>
                    )}
                </div>

                {/* Sección de Pagos */}
                <div className="dashboard-section pagos-section">
                    <div className="section-header">
                        <h2>Pagos Recientes</h2>
                        <button
                            className="btn-ver-mas"
                            onClick={() => navigate('/pagos')}
                        >
                            Ver historial →
                        </button>
                    </div>

                    {pagosRecientes.length > 0 ? (
                        <div className="pagos-lista-mini">
                            {pagosRecientes.map((pago) => (
                                <div key={pago.id} className="pago-mini">
                                    <div className="pago-mini-info">
                                        <span className="pago-icono">
                                            {pago.payment_method === 'credit_card' ? <CreditCard size={20} /> : <DollarSign size={20} />}
                                        </span>
                                        <div>
                                            <span className="pago-concepto">
                                                {pago.metadata?.plan_nombre || 'Pago'}
                                            </span>
                                            <span className="pago-fecha">
                                                {new Date(pago.created_at).toLocaleDateString('es-AR')}
                                            </span>
                                        </div>
                                    </div>
                                    <div className="pago-mini-monto">
                                        <span className="monto">${pago.amount?.toFixed(2)}</span>
                                        <span className={`estado-pago ${pago.status}`}>
                                            {pago.status === 'completed' ? <CheckCircle2 size={16} /> : <Clock size={16} />}
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div className="no-pagos-widget">
                            <p>No hay pagos registrados</p>
                        </div>
                    )}
                </div>
            </div>

            <div className="dashboard-acciones-rapidas">
                <h2>Acciones Rápidas</h2>
                <div className="acciones-grid">
                    <button
                        className="accion-card"
                        onClick={() => navigate('/planes')}
                    >
                        <span className="accion-icono"><ClipboardList size={24} /></span>
                        <span className="accion-texto">Ver Planes</span>
                    </button>
                    <button
                        className="accion-card"
                        onClick={() => navigate('/actividades')}
                    >
                        <span className="accion-icono"><Dumbbell size={24} /></span>
                        <span className="accion-texto">Actividades</span>
                    </button>
                    <button
                        className="accion-card"
                        onClick={() => navigate('/sucursales')}
                    >
                        <span className="accion-icono"><MapPin size={24} /></span>
                        <span className="accion-texto">Sucursales</span>
                    </button>
                    <button
                        className="accion-card"
                        onClick={() => navigate('/pagos')}
                    >
                        <span className="accion-icono"><CreditCard size={24} /></span>
                        <span className="accion-texto">Mis Pagos</span>
                    </button>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
