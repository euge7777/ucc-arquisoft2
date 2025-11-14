import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ClipboardList, Check, AlertCircle } from 'lucide-react';
import { getMockSuscripcionByUserId, mockApiDelay } from '../data/mockData';
import '../styles/MiSuscripcion.css';

const MiSuscripcion = () => {
    const [suscripcion, setSuscripcion] = useState(null);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();
    const userId = localStorage.getItem("idUsuario");

    useEffect(() => {
        const fetchSuscripcion = async () => {
            try {
                await mockApiDelay();
                const sub = getMockSuscripcionByUserId(userId);
                setSuscripcion(sub);
            } catch (error) {
                console.error("Error al cargar suscripción:", error);
            } finally {
                setLoading(false);
            }
        };

        fetchSuscripcion();
    }, [userId]);

    const handleCancelar = async () => {
        if (!window.confirm("¿Estás seguro de que deseas cancelar tu suscripción?")) {
            return;
        }

        try {
            await mockApiDelay();
            // Aquí iría la llamada real a la API cuando esté lista
            // await fetch(SUBSCRIPTIONS_API.subscriptionById(suscripcion.id), { method: 'DELETE' })
            alert("Suscripción cancelada exitosamente");
            setSuscripcion({ ...suscripcion, estado: "cancelada" });
        } catch (error) {
            console.error("Error al cancelar suscripción:", error);
            alert("Error al cancelar la suscripción");
        }
    };

    const handleRenovar = () => {
        if (suscripcion && suscripcion.plan) {
            navigate(`/checkout/${suscripcion.plan.id}`);
        }
    };

    const getDiasRestantes = (fechaVencimiento) => {
        const hoy = new Date();
        const vencimiento = new Date(fechaVencimiento);
        const diffTime = vencimiento - hoy;
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        return diffDays;
    };

    const getEstadoBadgeClass = (estado) => {
        switch (estado) {
            case 'activa':
                return 'estado-activa';
            case 'vencida':
                return 'estado-vencida';
            case 'cancelada':
                return 'estado-cancelada';
            case 'pendiente_pago':
                return 'estado-pendiente';
            default:
                return '';
        }
    };

    if (loading) {
        return (
            <div className="mi-suscripcion-container">
                <div className="loading-message">Cargando suscripción...</div>
            </div>
        );
    }

    if (!suscripcion) {
        return (
            <div className="mi-suscripcion-container">
                <div className="no-suscripcion">
                    <div className="no-suscripcion-icon">
                        <ClipboardList size={48} />
                    </div>
                    <h2>No tenés una suscripción activa</h2>
                    <p>Suscribite a uno de nuestros planes para acceder a todas las actividades</p>
                    <button
                        className="btn-ver-planes"
                        onClick={() => navigate('/planes')}
                    >
                        Ver Planes Disponibles
                    </button>
                </div>
            </div>
        );
    }

    const diasRestantes = getDiasRestantes(suscripcion.fecha_vencimiento);
    const proximoVencer = diasRestantes <= 7 && diasRestantes > 0;

    return (
        <div className="mi-suscripcion-container">
            <div className="suscripcion-header">
                <h1>Mi Suscripción</h1>
            </div>

            <div className="suscripcion-content">
                <div className="suscripcion-card">
                    <div className="suscripcion-plan-header" style={{ borderTopColor: suscripcion.plan?.color }}>
                        <div className="plan-info">
                            <h2>{suscripcion.plan?.nombre || "Plan"}</h2>
                            <p className="plan-descripcion">{suscripcion.plan?.descripcion}</p>
                        </div>
                        <div className={`estado-badge ${getEstadoBadgeClass(suscripcion.estado)}`}>
                            {suscripcion.estado.toUpperCase()}
                        </div>
                    </div>

                    <div className="suscripcion-detalles">
                        <div className="detalle-item">
                            <span className="detalle-label">Precio Mensual:</span>
                            <span className="detalle-valor precio">${suscripcion.plan?.precio_mensual.toFixed(2)}</span>
                        </div>
                        <div className="detalle-item">
                            <span className="detalle-label">Fecha de Inicio:</span>
                            <span className="detalle-valor">
                                {new Date(suscripcion.fecha_inicio).toLocaleDateString('es-AR')}
                            </span>
                        </div>
                        <div className="detalle-item">
                            <span className="detalle-label">Fecha de Vencimiento:</span>
                            <span className="detalle-valor">
                                {new Date(suscripcion.fecha_vencimiento).toLocaleDateString('es-AR')}
                            </span>
                        </div>
                        <div className="detalle-item">
                            <span className="detalle-label">Días Restantes:</span>
                            <span className={`detalle-valor ${proximoVencer ? 'dias-advertencia' : ''}`}>
                                {diasRestantes > 0 ? `${diasRestantes} días` : 'Vencida'}
                            </span>
                        </div>
                        {suscripcion.metadata?.auto_renovacion && (
                            <div className="detalle-item">
                                <span className="detalle-label">Auto-renovación:</span>
                                <span className="detalle-valor renovacion-activa"><Check size={16} className="inline mr-1" /> Activada</span>
                            </div>
                        )}
                    </div>

                    {proximoVencer && (
                        <div className="alerta-vencimiento">
                            <AlertCircle size={20} className="inline mr-2" />
                            Tu suscripción vence pronto. Renovála para seguir disfrutando de los beneficios.
                        </div>
                    )}

                    <div className="suscripcion-beneficios">
                        <h3>Beneficios de tu plan:</h3>
                        <ul>
                            {suscripcion.plan?.beneficios?.map((beneficio, index) => (
                                <li key={index}>
                                    <span className="check-icon"><Check size={20} /></span>
                                    {beneficio}
                                </li>
                            ))}
                        </ul>
                    </div>

                    <div className="suscripcion-acciones">
                        {suscripcion.estado === 'activa' && (
                            <>
                                <button className="btn-renovar" onClick={handleRenovar}>
                                    Renovar Suscripción
                                </button>
                                <button className="btn-cancelar" onClick={handleCancelar}>
                                    Cancelar Suscripción
                                </button>
                            </>
                        )}
                        {(suscripcion.estado === 'vencida' || suscripcion.estado === 'cancelada') && (
                            <button className="btn-renovar-principal" onClick={() => navigate('/planes')}>
                                Ver Planes Disponibles
                            </button>
                        )}
                        {suscripcion.estado === 'pendiente_pago' && (
                            <button className="btn-pagar" onClick={() => navigate('/pagos')}>
                                Completar Pago
                            </button>
                        )}
                    </div>
                </div>

                {suscripcion.historial_renovaciones && suscripcion.historial_renovaciones.length > 0 && (
                    <div className="historial-renovaciones">
                        <h3>Historial de Renovaciones</h3>
                        <div className="renovaciones-lista">
                            {suscripcion.historial_renovaciones.map((renovacion, index) => (
                                <div key={index} className="renovacion-item">
                                    <div className="renovacion-fecha">
                                        {new Date(renovacion.fecha).toLocaleDateString('es-AR')}
                                    </div>
                                    <div className="renovacion-monto">${renovacion.monto.toFixed(2)}</div>
                                    <div className="renovacion-pago">ID: {renovacion.pago_id}</div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default MiSuscripcion;
