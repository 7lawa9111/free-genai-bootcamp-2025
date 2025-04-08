FROM opea/speecht5:latest

USER root

# Uninstall jax and jaxlib completely
RUN pip uninstall -y jax jaxlib || true

# Set environment variables to make docarray skip JAX
ENV DOCARRAY_SKIP_JAX=1 \
    DOCARRAY_COMPUTE_BACKEND=numpy \
    DOCARRAY_NO_WARNINGS=1

# Switch back to non-root user if needed
USER 1000 